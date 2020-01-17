package process

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"github.com/google/uuid"
)

type (
	// Provider for Processes
	Provider func() *Process

	// Process representing a place order process and has a current context with infos about result and current state
	Process struct {
		context     Context
		allStates   map[string]State
		failedState FailedState
		logger      flamingo.Logger
	}

	// Factory use to get Process instance
	Factory struct {
		provider    Provider
		startState  State
		failedState FailedState
	}

	// RollbackReference a reference that can be used to trigger a rollback
	RollbackReference struct {
		StateName string
		Data      RollbackData
	}

	// RollbackData needed for rollback of a state
	RollbackData interface{}

	// FailedReason gives a human readable reason for a state failure
	FailedReason interface {
		Reason() string
	}

	// ErrorOccurredReason is used for unspecified errors
	ErrorOccurredReason struct {
		Error string
	}

	// PaymentErrorOccurredReason is used for errors during payment
	PaymentErrorOccurredReason struct {
		Error string
	}

	// CartValidationErrorReason contains the ValidationResult
	CartValidationErrorReason struct {
		ValidationResult validation.Result
	}
)

func init() {
	gob.Register(ErrorOccurredReason{})
	gob.Register(PaymentErrorOccurredReason{})
	gob.Register(CartValidationErrorReason{})
}

// Reason for the error occurred
func (e ErrorOccurredReason) Reason() string {
	return e.Error
}

// Reason for the error occurred
func (e PaymentErrorOccurredReason) Reason() string {
	return e.Error
}

// Reason for failing
func (e CartValidationErrorReason) Reason() string {
	return "Cart invalid"
}

// Inject dependencies
func (f *Factory) Inject(
	provider Provider,
	failedState FailedState,
	dep *struct {
		StartState State `inject:"startState"`
	},
) {
	f.provider = provider
	f.failedState = failedState
	if dep != nil {
		f.startState = dep.StartState
	}
}

// New process with initial state
func (f *Factory) New(returnURL *url.URL) (*Process, error) {
	if f.startState == nil {
		return nil, errors.New("no start state given")
	}
	p := f.provider()
	p.failedState = f.failedState
	p.context = Context{
		UUID:      uuid.New().String(),
		State:     f.startState,
		ReturnURL: returnURL,
	}

	return p, nil
}

// NewFromProcessContext returns a new process with given Context
func (f *Factory) NewFromProcessContext(pctx Context) (*Process, error) {
	p := f.provider()
	p.failedState = f.failedState
	p.context = pctx

	return p, nil
}

// Inject dependencies
func (p *Process) Inject(
	allStates map[string]State,
	logger flamingo.Logger,
) *Process {
	p.allStates = allStates
	p.logger = logger.
		WithField(flamingo.LogKeyModule, "checkout").
		WithField(flamingo.LogKeyCategory, "process")

	return p
}

// Run triggers run on current state
func (p *Process) Run(ctx context.Context) {
	runResult := p.context.State.Run(ctx, p)
	if runResult.RollbackData != nil {
		p.context.RollbackReferences = append(p.context.RollbackReferences, RollbackReference{
			StateName: p.context.State.Name(),
			Data:      runResult.RollbackData,
		})
	}
	if runResult.Failed != nil {
		p.Failed(ctx, runResult.Failed)
	}
}

func (p *Process) rollback() error {
	for i := len(p.context.RollbackReferences) - 1; i >= 0; i-- {
		rollbackRef := p.context.RollbackReferences[i]
		state, ok := p.allStates[rollbackRef.StateName]
		if !ok {
			p.logger.Error(fmt.Errorf("state %q not found for rollback", rollbackRef.StateName))
			continue
		}
		// todo error types for fatal end and continue rollback chain
		_ = state.Rollback(rollbackRef.Data)
	}

	return nil
}

// Context to get current process context
func (p *Process) Context() Context {
	return p.context
}

// UpdateState updates
func (p *Process) UpdateState(s State) {
	p.context.State = s
}

// Failed performs all collected rollbacks and switches to FailedState
func (p *Process) Failed(ctx context.Context, reason FailedReason) {
	err := p.rollback()
	if err != nil {
		p.logger.WithContext(ctx).Error("rollback failed: ", err)
	}
	p.UpdateState(p.failedState.SetFailedReason(reason))
}
