package process

import (
	"context"
)

type (
	// State interface
	State interface {
		Run(context.Context, *Process, StateData) RunResult
		Rollback(context.Context, RollbackData) error
		IsFinal() bool
		Name() string
	}

	// RunResult of a state
	RunResult struct {
		RollbackData RollbackData
		Failed       FailedReason
	}
)
