package bidding_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"github.com/students-api/bidding-service/internal/services/bidding_service"
	"github.com/students-api/bidding-service/internal/services/bidding_service/mocks"
)

func TestPlaceBid_Success(t *testing.T) {
	mockRepo := new(mocks.BidRepository)
	mockProducer := new(mocks.EventSender)
	mockLock := new(mocks.LockService)

	mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(true, nil)
	mockLock.On("ReleaseLock", mock.Anything, int64(100)).Return(nil)
	mockRepo.On("CreateBid", mock.Anything, mock.MatchedBy(func(b interface{}) bool {
		return true 
	})).Return(nil)
	mockProducer.On("SendBidPlaced", mock.Anything).Return(nil)

	service := bidding_service.NewBiddingService(mockRepo, mockProducer, mockLock)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, UserId: 1, Amount: 1000}
	resp, err := service.PlaceBid(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Ставка принята", resp.Message)

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestPlaceBid_InvalidAmount(t *testing.T) {
	mockRepo := new(mocks.BidRepository)
	mockProducer := new(mocks.EventSender)
	mockLock := new(mocks.LockService)

	service := bidding_service.NewBiddingService(mockRepo, mockProducer, mockLock)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: -500}
	resp, err := service.PlaceBid(context.Background(), req)

	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Message, "больше нуля")

	mockRepo.AssertNotCalled(t, "CreateBid")
}

func TestPlaceBid_Locked(t *testing.T) {
	mockRepo := new(mocks.BidRepository)
	mockProducer := new(mocks.EventSender)
	mockLock := new(mocks.LockService)

	mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(false, nil)

	service := bidding_service.NewBiddingService(mockRepo, mockProducer, mockLock)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: 1000}
	resp, err := service.PlaceBid(context.Background(), req)

	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Message, "Сервер перегружен")
}

func TestPlaceBid_DBError(t *testing.T) {
	mockRepo := new(mocks.BidRepository)
	mockProducer := new(mocks.EventSender)
	mockLock := new(mocks.LockService)

	mockLock.On("AcquireLock", mock.Anything, int64(100)).Return(true, nil)
	mockLock.On("ReleaseLock", mock.Anything, int64(100)).Return(nil)
	
	mockRepo.On("CreateBid", mock.Anything, mock.Anything).Return(errors.New("db error"))

	service := bidding_service.NewBiddingService(mockRepo, mockProducer, mockLock)

	req := &bidding_api.PlaceBidRequest{AuctionId: 100, Amount: 1000}
	resp, err := service.PlaceBid(context.Background(), req)

	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Message, "Ошибка БД")
}