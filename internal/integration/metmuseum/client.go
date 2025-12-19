package metmuseum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://collectionapi.metmuseum.org/public/collection/v1"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ObjectResponse struct {
	ObjectID     int64  `json:"objectID"`
	Title        string `json:"title"`
	Artist       string `json:"artistDisplayName"`
	PrimaryImage string `json:"primaryImage"`
}

func (c *Client) GetObjectData(objectID int64) (*ObjectResponse, error) {
	url := fmt.Sprintf("%s/objects/%d", BaseURL, objectID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("museum api returned status: %d", resp.StatusCode)
	}

	var data ObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}