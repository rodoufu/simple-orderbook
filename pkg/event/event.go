package event

// Event represents an occurrence that needs to be propagated.
type Event interface {
	Output
	event()
}
