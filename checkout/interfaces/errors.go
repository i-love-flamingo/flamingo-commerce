package interfaces

import "errors"

var (
	ErrCheckoutGeneral = errors.New("checkout_general_error")

	// ErrNoPlaceOrderProcess is returned when no place order process exists.
	// Keep the error message stable for existing GraphQL clients.
	ErrNoPlaceOrderProcess = errors.New("ErrNoPlaceOrderProcess")

	// ErrCancelNotPossibleFinalState is returned when a cancel is attempted for a final state process.
	// Keep the error message stable for existing GraphQL clients.
	ErrCancelNotPossibleFinalState = errors.New("process already in final state, cancel not possible")
)
