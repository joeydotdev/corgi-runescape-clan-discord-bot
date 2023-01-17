package handlers

// Handler is a struct that contains custom metadata around a Discord Event.
type Handler struct{}

// New creates a new Handler.
func New() *Handler {
	return &Handler{}
}
