package domain

// CancellationReason describes why a payment was cancelled. It is optional at
// every layer; the zero value (Unspecified) preserves legacy behaviour.
type CancellationReason int

const (
	// CancellationReasonUnspecified is the default / legacy behaviour.
	CancellationReasonUnspecified CancellationReason = iota
	// CancellationReasonAbortedByCustomer marks an explicit customer abort.
	CancellationReasonAbortedByCustomer
)
