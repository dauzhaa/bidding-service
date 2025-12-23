package bidding_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"github.com/students-api/bidding-service/internal/services/bidding_service"
	"github.com/students-api/bidding-service/internal/services/bidding_service/mocks"
)

type BiddingServiceSuite struct {
	suite.Suite
	mockRepo     *mocks.BidRepository
	mockProducer *mocks.EventSender
	mockLock     *mocks.LockService
	service      bidding_api.BiddingServiceServer
}

func (s *BiddingServiceSuite) SetupTest() {
	s.mockRepo = new(mocks.BidRepository)
	s.mockProducer = new(mocks.EventSender)
	s.mockLock = new(mocks.LockService)

	s.service = bidding_service.NewBiddingService(s.mockRepo, s.mockProducer, s.mockLock)
}

func (s *BiddingServiceSuite) TestPlaceBid_Success() {
	s.mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(true, nil)
	s.mockLock.On("ReleaseLock", mock.Anything, int64(100)).Return(nil)
	s.mockRepo.On("CreateBid", mock.Anything, mock.MatchedBy(func(b interface{}) bool {
		return true
	})).Return(nil)
	s.mockProducer.On("SendBidPlaced", mock.Anything).Return(nil)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, UserId: 1, Amount: 1000}
	resp, err := s.service.PlaceBid(context.Background(), req)

	s.NoError(err)
	s.True(resp.Success)
	s.Equal("Ставка принята", resp.Message)

	s.mockRepo.AssertExpectations(s.T())
	s.mockProducer.AssertExpectations(s.T())
}

func (s *BiddingServiceSuite) TestPlaceBid_InvalidAmount() {
	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: -500}
	resp, err := s.service.PlaceBid(context.Background(), req)

	s.NoError(err)
	s.False(resp.Success)
	s.Contains(resp.Message, "больше нуля")

	s.mockRepo.AssertNotCalled(s.T(), "CreateBid")
}

func (s *BiddingServiceSuite) TestPlaceBid_Locked() {
	s.mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(false, nil)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: 1000}
	resp, err := s.service.PlaceBid(context.Background(), req)

	s.NoError(err)
	s.False(resp.Success)
	s.Contains(resp.Message, "Сервер перегружен")
}

func (s *BiddingServiceSuite) TestPlaceBid_DBError() {
	s.mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(true, nil)
	s.mockLock.On("ReleaseLock", mock.Anything, int64(100)).Return(nil)
	
	s.mockRepo.On("CreateBid", mock.Anything, mock.Anything).Return(errors.New("db error"))

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: 1000}
	resp, err := s.service.PlaceBid(context.Background(), req)

	s.NoError(err)
	s.False(resp.Success)
	s.Contains(resp.Message, "Ошибка БД")
}

func TestBiddingServiceSuite(t *testing.T) {
	suite.Run(t, new(BiddingServiceSuite))
}