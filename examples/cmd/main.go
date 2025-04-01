package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/stremovskyy/pulser"
)

func main() {
	logger := log.New(os.Stdout, "[PULSE] ", log.LstdFlags)

	config := pulser.Config{
		APIURL:      "http://localhost:8080/pulse/api/v1/srv/track",
		APIKey:      os.Getenv("PULSE_API_KEY"),
		ServiceCode: "mobile-app",
		Timeout:     5 * time.Second,
		MaxRetries:  2,
		LogFunc: func(format string, args ...interface{}) {
			logger.Printf(format, args...)
		},
		AsyncEnabled: true,
		QueueSize:    1000,
	}

	client, err := pulser.NewClient(config)
	if err != nil {
		logger.Fatalf("Failed to create metrics client: %v", err)
	}
	defer func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Shutdown(5 * time.Second); err != nil {
			logger.Printf("Error during metrics client shutdown: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.Track(ctx, "user123", pulser.EventTypePushOpened, pulser.EventSubTypeProduction, nil)
	if err != nil {
		logger.Printf("Failed to track push open: %v", err)
	}

	metadata := map[string]interface{}{
		"app_version":    "1.2.3",
		"os_version":     "iOS 15.4",
		"device_model":   "iPhone 13",
		"network_type":   "wifi",
		"session_id":     "abc-123-def-456",
		"time_spent_sec": 45,
	}

	err = client.Track(ctx, "user456", pulser.EventTypeAppOpened, pulser.EventSubTypeAnalytics, metadata)
	if err != nil {
		logger.Printf("Failed to track app open: %v", err)
	}

	formMetadata := map[string]interface{}{
		"form_id":              "checkout",
		"form_step":            3,
		"completed":            true,
		"time_to_complete_sec": 120,
		"validation_errors":    0,
	}

	err = client.Track(ctx, "user789", pulser.EventTypeFormSubmit, pulser.EventSubTypeEngagement, formMetadata)
	if err != nil {
		logger.Printf("Failed to track form submit: %v", err)
	}

	errorMetadata := map[string]interface{}{
		"error_code":    "E1001",
		"error_message": "Failed to connect to payment gateway",
		"component":     "payment-service",
		"severity":      "high",
		"retry_count":   2,
	}

	err = client.Track(ctx, "user101", pulser.EventTypeError, pulser.EventSubTypeProduction, errorMetadata)
	if err != nil {
		logger.Printf("Failed to track error: %v", err)
	}

	trackScreenView := func(userID, screenName string, params map[string]interface{}) {
		metadata := map[string]interface{}{
			"screen_name": screenName,
			"timestamp":   time.Now().Unix(),
		}

		for k, v := range params {
			metadata[k] = v
		}

		if err := client.Track(ctx, userID, pulser.EventTypeScreenView, pulser.EventSubTypeAnalytics, metadata); err != nil {
			logger.Printf("Failed to track screen view: %v", err)
		}
	}

	trackScreenView(
		"user202", "product_details", map[string]interface{}{
			"product_id": "prod-123",
			"category":   "electronics",
			"referrer":   "search_results",
		},
	)

	time.Sleep(1 * time.Second)
}
