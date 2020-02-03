package process

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/payment/application"
)

type (
	// State interface
	State interface {
		Run(context.Context, *Process) RunResult
		Rollback(context.Context, RollbackData) error
		IsFinal() bool
		Name() string
	}

	// RunResult of a state
	RunResult struct {
		RollbackData RollbackData
		Failed       FailedReason
	}

	// FatalRollbackError which causes the premature end of rollback process
	FatalRollbackError struct {
		error error
		State string
	}

	// PaymentValidatorFunc to decide over next state depending on payment situation
	PaymentValidatorFunc func(ctx context.Context, p *Process, paymentService *application.PaymentService) RunResult
)

// NewFatalRollbackError mark a rollback error as fatal, this breaks the current rollback
func NewFatalRollbackError(err error, stateName string) FatalRollbackError {
	return FatalRollbackError{
		error: err,
		State: stateName,
	}
}

func (f *FatalRollbackError) Error() string {
	return f.error.Error()
}
