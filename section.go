package ptvvisum

// Section interface for all section types
type Section interface {
	Name() string
	Headers() []string
}
