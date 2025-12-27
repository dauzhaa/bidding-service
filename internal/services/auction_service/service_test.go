package auction_service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/students-api/bidding-service/internal/pb/auction_api"
)

type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) CreateAuction(ctx context.Context, auction *auction_api.Auction) error {
	args := m.Called(ctx, auction)
	return args.Error(0)
}

func (m *MockAuctionRepository) ListAuctions(ctx context.Context) ([]*auction_api.Auction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auction_api.Auction), args.Error(1)
}


func TestListAuctions_Success(t *testing.T) {
	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	ctx := context.Background()

	expectedAuctions := []*auction_api.Auction{
		{Id: 1, Title: "Test Art", CurrentPrice: 100},
		{Id: 2, Title: "Test Art 2", CurrentPrice: 200},
	}

	mockRepo.On("ListAuctions", ctx).Return(expectedAuctions, nil)

	resp, err := service.ListAuctions(ctx, &auction_api.ListAuctionsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Auctions))
	mockRepo.AssertExpectations(t)
}

func TestListAuctions_RepoError(t *testing.T) {
	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	ctx := context.Background()

	mockRepo.On("ListAuctions", ctx).Return(nil, errors.New("database down"))

	resp, err := service.ListAuctions(ctx, &auction_api.ListAuctionsRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database down")
	mockRepo.AssertExpectations(t)
}

func TestCreateAuction_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"objectID": 436535,
			"title": "Wheat Field",
			"artistDisplayName": "Vincent van Gogh",
			"primaryImage": "http://example.com/image.jpg"
		}`))
	}))
	defer mockServer.Close()

	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	service.APIURL = mockServer.URL

	ctx := context.Background()
	req := &auction_api.CreateAuctionRequest{
		ObjectId:   436535,
		StartPrice: 1000,
	}

	mockRepo.On("CreateAuction", ctx, mock.AnythingOfType("*auction_api.Auction")).Return(nil)

	response, err := service.CreateAuction(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Wheat Field", response.Title)
	mockRepo.AssertExpectations(t)
}

func TestCreateAuction_ApiError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	service.APIURL = mockServer.URL

	req := &auction_api.CreateAuctionRequest{ObjectId: 999}
	response, err := service.CreateAuction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "Met API returned status: 404")
}

func TestCreateAuction_RepoError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"objectID": 1, "title": "Art", "artistDisplayName": "Me"}`))
	}))
	defer mockServer.Close()

	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	service.APIURL = mockServer.URL

	mockRepo.On("CreateAuction", mock.Anything, mock.Anything).Return(errors.New("insert failed"))

	req := &auction_api.CreateAuctionRequest{ObjectId: 1}
	response, err := service.CreateAuction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "insert failed")
}

func TestCreateAuction_DefaultValues(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"objectID": 2, "title": "", "artistDisplayName": ""}`))
	}))
	defer mockServer.Close()

	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	service.APIURL = mockServer.URL

	mockRepo.On("CreateAuction", mock.Anything, mock.Anything).Return(nil)

	req := &auction_api.CreateAuctionRequest{ObjectId: 2}
	response, err := service.CreateAuction(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Unknown Title", response.Title)
	assert.Equal(t, "Unknown Artist", response.Artist)
}

func TestCreateAuction_InvalidJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid-json`))
	}))
	defer mockServer.Close()

	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(mockRepo)
	service.APIURL = mockServer.URL

	req := &auction_api.CreateAuctionRequest{ObjectId: 3}
	response, err := service.CreateAuction(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to decode")
}