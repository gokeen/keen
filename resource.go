package keen

import (
	"fmt"
)

// Resource represents the interface for a keen API resource
type Resource interface {
	Authorization() string
	Path() string
	Data() interface{}
}

var (
	resourcefmt = fmt.
		Sprintf("/%s/projects/%s/%s/%s", API_VERSION, "%s", "%s", "%s")
)

// QueryResource implements a struct for query based read calls
type QueryResource struct {
	Query
	Key       string
	ProjectID string
}

func (q QueryResource) Authorization() string {
	return q.Key
}

func (q QueryResource) Path() string {
	return fmt.
		Sprintf(resourcefmt, q.ProjectID, "queries", q.Query.QueryType())
}

func (q QueryResource) Data() interface{} {
	return q.Query
}

// EventResource implements a struct for event based write calls
type EventResource struct {
	Event
	Key       string
	ProjectID string
}

func (q EventResource) Authorization() string {
	return q.Key
}

func (q EventResource) Path() string {
	return fmt.
		Sprintf(resourcefmt, q.ProjectID, "events", q.Event.CollectionName())
}

func (q EventResource) Data() interface{} {
	return q.Event
}
