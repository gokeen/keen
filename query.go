package keen

// Query represents the interface for a analysis query
type Query interface {
	QueryType() string
}
