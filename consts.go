package pulser

type EventType string

type EventSubType string

const (
	EventTypeAppOpened     EventType = "app_opened"
	EventTypeAppClosed     EventType = "app_closed"
	EventTypeAppBackground EventType = "app_background"
	EventTypeAppForeground EventType = "app_foreground"
	EventTypeAppCrash      EventType = "app_crash"

	EventTypePushOpened    EventType = "push_opened"
	EventTypePushReceived  EventType = "push_received"
	EventTypePushDismissed EventType = "push_dismissed"

	EventTypeUserLogin      EventType = "user_login"
	EventTypeUserLogout     EventType = "user_logout"
	EventTypeUserRegistered EventType = "user_registered"

	EventTypeScreenView    EventType = "screen_view"
	EventTypeButtonClick   EventType = "button_click"
	EventTypeFormSubmit    EventType = "form_submit"
	EventTypeItemSelection EventType = "item_selection"
	EventTypeSearch        EventType = "search"

	EventTypePaymentInitiated  EventType = "payment_initiated"
	EventTypePaymentCompleted  EventType = "payment_completed"
	EventTypePaymentFailed     EventType = "payment_failed"
	EventTypeCheckoutStarted   EventType = "checkout_started"
	EventTypeCheckoutCompleted EventType = "checkout_completed"
	EventTypeCheckoutAbandoned EventType = "checkout_abandoned"

	EventTypeFeatureEnabled  EventType = "feature_enabled"
	EventTypeFeatureDisabled EventType = "feature_disabled"
	EventTypeFeatureUsed     EventType = "feature_used"

	EventTypeError           EventType = "error"
	EventTypeAPIError        EventType = "api_error"
	EventTypeNetworkError    EventType = "network_error"
	EventTypeValidationError EventType = "validation_error"
)

const (
	EventSubTypeTest       EventSubType = "test"
	EventSubTypeProduction EventSubType = "production"
	EventSubTypeStaging    EventSubType = "staging"
	EventSubTypeDebug      EventSubType = "debug"

	EventSubTypeAnalytics   EventSubType = "analytics"
	EventSubTypePerformance EventSubType = "performance"
	EventSubTypeUsage       EventSubType = "usage"
	EventSubTypeBehavior    EventSubType = "behavior"

	EventSubTypeOnboarding EventSubType = "onboarding"
	EventSubTypeConversion EventSubType = "conversion"
	EventSubTypeEngagement EventSubType = "engagement"
	EventSubTypeRetention  EventSubType = "retention"

	EventSubTypeCore         EventSubType = "core"
	EventSubTypePremium      EventSubType = "premium"
	EventSubTypeBeta         EventSubType = "beta"
	EventSubTypeExperimental EventSubType = "experimental"
)
