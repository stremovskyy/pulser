package pulser

import (
	"context"
	"time"
)

type Client interface {
	Track(ctx context.Context, userID string, eventType EventType, eventSubType EventSubType, metadata map[string]interface{}) error
	TrackWithServiceCode(ctx context.Context, serviceCode, userID string, eventType EventType, eventSubType EventSubType, metadata map[string]interface{}) error
	Shutdown(timeout time.Duration) error
}
