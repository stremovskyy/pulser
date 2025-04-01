package pulser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Config struct {
	APIURL      string
	APIKey      string
	ServiceCode string

	// opts
	Timeout           time.Duration
	MaxRetries        int
	RetryWaitTime     time.Duration
	RetryMaxWaitTime  time.Duration
	LogFunc           func(string, ...interface{})
	AsyncEnabled      bool
	QueueSize         int
	DisableKeepAlives bool
}

type Event struct {
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ServiceCode  string                 `json:"service_code"`
	UserID       string                 `json:"user_id"`
	EventSubType EventSubType           `json:"event_sub_type"`
	EventType    EventType              `json:"event_type"`
}

type client struct {
	config     Config
	httpClient *http.Client

	queue       chan Event
	wg          sync.WaitGroup
	ctx         context.Context
	cancelFunc  context.CancelFunc
	initialized bool
	mu          sync.Mutex
}

type errResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func DefaultConfig() Config {
	return Config{
		APIURL:            "",
		Timeout:           10 * time.Second,
		MaxRetries:        3,
		RetryWaitTime:     100 * time.Millisecond,
		RetryMaxWaitTime:  2 * time.Second,
		LogFunc:           nil,
		AsyncEnabled:      false,
		QueueSize:         1000,
		DisableKeepAlives: false,
	}
}

func NewClient(config Config) (Client, error) {
	if config.APIURL == "" {
		return nil, errors.New("APIURL is required")
	}
	if config.APIKey == "" {
		return nil, errors.New("APIKey is required")
	}
	if config.ServiceCode == "" {
		return nil, errors.New("ServiceCode is required")
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultConfig().Timeout
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = DefaultConfig().MaxRetries
	}
	if config.RetryWaitTime == 0 {
		config.RetryWaitTime = DefaultConfig().RetryWaitTime
	}
	if config.RetryMaxWaitTime == 0 {
		config.RetryMaxWaitTime = DefaultConfig().RetryMaxWaitTime
	}
	if config.QueueSize == 0 {
		config.QueueSize = DefaultConfig().QueueSize
	}

	_, err := url.Parse(config.APIURL)
	if err != nil {
		return nil, fmt.Errorf("invalid APIURL: %w", err)
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives: config.DisableKeepAlives,
		},
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	client := &client{
		config:      config,
		httpClient:  httpClient,
		ctx:         ctx,
		cancelFunc:  cancelFunc,
		initialized: true,
	}

	if config.AsyncEnabled {
		client.queue = make(chan Event, config.QueueSize)
		client.startWorker()
	}

	return client, nil
}

func (c *client) log(format string, args ...interface{}) {
	if c.config.LogFunc != nil {
		c.config.LogFunc(format, args...)
	}
}

func (c *client) Track(ctx context.Context, userID string, eventType EventType, eventSubType EventSubType, metadata map[string]interface{}) error {
	return c.TrackWithServiceCode(ctx, c.config.ServiceCode, userID, eventType, eventSubType, metadata)
}
func (c *client) TrackWithServiceCode(ctx context.Context, serviceCode, userID string, eventType EventType, eventSubType EventSubType, metadata map[string]interface{}) error {
	if !c.initialized {
		return errors.New("client not properly initialized")
	}

	event := Event{
		ServiceCode:  serviceCode,
		UserID:       userID,
		EventType:    eventType,
		EventSubType: eventSubType,
		Metadata:     metadata,
	}

	if c.config.AsyncEnabled {
		select {
		case c.queue <- event:
			return nil
		default:
			c.log("Metrics queue full, dropping event: %+v", event)
			return errors.New("async queue full, event dropped")
		}
	}

	return c.sendEvent(ctx, event)
}

func (c *client) sendEvent(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.APIURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-API-KEY", c.config.APIKey)

	c.log("Sending metric event: %s", string(payload))

	var lastErr error
	for retries := 0; retries <= c.config.MaxRetries; retries++ {
		if retries > 0 {
			waitTime := c.config.RetryWaitTime * time.Duration(1<<uint(retries-1))
			if waitTime > c.config.RetryMaxWaitTime {
				waitTime = c.config.RetryMaxWaitTime
			}

			c.log("Retrying request after %v (attempt %d/%d)", waitTime, retries, c.config.MaxRetries)

			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			c.log("Request error: %v", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			c.log("Metric event sent successfully")
			return nil
		}

		respBody, _ := io.ReadAll(resp.Body)

		var errResp errResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
			lastErr = fmt.Errorf("API error (%d): %s (code: %s)", resp.StatusCode, errResp.Message, errResp.Code)
		} else {
			lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
		}

		c.log("API error: %v", lastErr)

		if resp.StatusCode == http.StatusUnauthorized ||
			resp.StatusCode == http.StatusForbidden ||
			resp.StatusCode == http.StatusNotFound {
			return lastErr
		}
	}

	return fmt.Errorf("failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

func (c *client) startWorker() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.ctx.Done():
				c.log("Worker shutting down")
				return
			case event := <-c.queue:
				ctx, cancel := context.WithTimeout(c.ctx, c.config.Timeout)
				err := c.sendEvent(ctx, event)
				if err != nil {
					c.log("Error sending async event: %v", err)
				}
				cancel()
			}
		}
	}()
}

func (c *client) Shutdown(timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized || c.cancelFunc == nil {
		return nil
	}

	if !c.config.AsyncEnabled || c.queue == nil {
		c.cancelFunc()
		return nil
	}

	c.cancelFunc()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return errors.New("shutdown timed out waiting for events to process")
	}
}
