package models

// Session represents a tmux session.
type Session struct {
	// Name is the name of the session.
	Name string

	// Attached is whether the session is currently attached to a terminal.
	Attached bool

	// LastAttachedTime is the Unix timestamp of when the session was last attached to a terminal.
	LastAttachedTime int64
}
