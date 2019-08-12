package models

type FlowContext struct {
	TrackingID string
}

type Response struct {
	Data  interface{}
	Error *APIError
}

type APIError struct {
	Code     string
	Message  string
	Metadata interface{}
}
