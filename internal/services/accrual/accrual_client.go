package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akarashov/gofermart/internal/models"
)

var (
	ErrRateLimites = errors.New("rate limited")
	ErrServerError = errors.New("server error")
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	retryDelay time.Duration
	maxRetries int
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		retryDelay: 5 * time.Second,
		maxRetries: 10,
	}
}

// Get order info from accrual with retry and rate limit (svc)
func (c *Client) GetOrderInfo(ctx context.Context, orderNumber string) (*models.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)
	var resp *models.AccrualResponse
	var err error
	for i := 0; i < c.maxRetries; i++ {
		resp, err = c.doRequest(ctx, url)
		switch err {
		case nil:
			return resp, nil
		case ErrRateLimites:
			log.Println("Rate limited, retry")
			time.Sleep(c.retryDelay * time.Duration(i+1))
		case ErrServerError:
			log.Println("Accrual server error")
			return nil, err
		default:
			log.Println("Accrual unknown error")
			return nil, err
		}
	}
	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, err)
}

// HTTP Request to accrual service (repos)
func (c *Client) doRequest(ctx context.Context, url string) (*models.AccrualResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	// log.Printf("Get order data from accrual for new/processing orders %s\n", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var accrualResp models.AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
			return nil, err
		}
		return &accrualResp, nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		return nil, ErrRateLimites
	case http.StatusInternalServerError:
		return nil, ErrServerError
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}
