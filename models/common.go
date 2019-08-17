package models

// FlowContext use for flow level context
type FlowContext struct {
	TrackingID string
}

// Response generic API response
type Response struct {
	Data  interface{}
	Error *APIError
}

// APIError generic API error response
type APIError struct {
	Code     string
	Message  string
	Metadata interface{}
}
