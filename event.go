package keen

// Event represents the interface for an event object to be written to keen
type Event interface {
	ProjectID() string
	CollectionName() string
}

// EventResponse is the basic response structure for POST requests to
// https://keen.io/docs/api/reference/#event-collection-resource
type EventResponse struct {
	Created   bool   `json:"created"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}
