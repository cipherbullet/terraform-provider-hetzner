package types

// Client defines the interface that all API clients must implement
type Client interface {
	DoRequest(method, path string, payload interface{}) ([]byte, error)
	DoFormRequest(method, path string, values map[string][]string) ([]byte, error)
}
