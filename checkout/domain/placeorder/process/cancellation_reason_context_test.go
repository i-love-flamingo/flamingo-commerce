package process_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	paymentdomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

func TestCancellationReasonContext_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("absent key defaults to Unspecified", func(t *testing.T) {
		t.Parallel()

		got := process.CancellationReasonFromContext(context.Background())
		assert.Equal(t, paymentdomain.CancellationReasonUnspecified, got)
	})

	t.Run("value round-trips", func(t *testing.T) {
		t.Parallel()

		ctx := process.ContextWithCancellationReason(context.Background(), paymentdomain.CancellationReasonAbortedByCustomer)
		got := process.CancellationReasonFromContext(ctx)
		assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, got)
	})
}

func TestProcessContext_GobRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("with CancelReason set", func(t *testing.T) {
		t.Parallel()

		in := process.Context{UUID: "abc", CancelReason: paymentdomain.CancellationReasonAbortedByCustomer}

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(in))

		var out process.Context
		require.NoError(t, gob.NewDecoder(&buf).Decode(&out))
		assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, out.CancelReason)
	})

	t.Run("legacy bytes without reason decode to Unspecified", func(t *testing.T) {
		t.Parallel()

		in := process.Context{UUID: "abc"} // CancelReason left zero

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(in))

		var out process.Context
		require.NoError(t, gob.NewDecoder(&buf).Decode(&out))
		assert.Equal(t, paymentdomain.CancellationReasonUnspecified, out.CancelReason)
	})
}

// recordingState captures the reason it sees during rollback.
type recordingState struct {
	gotReason  paymentdomain.CancellationReason
	rolledBack bool
}

func (s *recordingState) Name() string { return "recording" }
func (s *recordingState) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{}
}
func (s *recordingState) IsFinal() bool { return false }
func (s *recordingState) Rollback(ctx context.Context, _ process.RollbackData) error {
	s.gotReason = process.CancellationReasonFromContext(ctx)
	s.rolledBack = true

	return nil
}

// stubFinalState is a minimal failed state so Process.Failed can switch to it.
type stubFinalState struct{}

func (stubFinalState) Name() string { return "Failed" }
func (stubFinalState) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{}
}
func (stubFinalState) IsFinal() bool                                        { return true }
func (stubFinalState) Rollback(context.Context, process.RollbackData) error { return nil }

func TestProcess_FailedWithReason_FlowsToRollback(t *testing.T) {
	t.Parallel()

	rec := &recordingState{}

	factory := &process.Factory{}
	factory.Inject(
		func() *process.Process {
			return (&process.Process{}).Inject(
				map[string]process.State{"recording": rec, "Failed": stubFinalState{}},
				flamingo.NullLogger{},
				nil,
			)
		},
		&struct {
			StartState  process.State `inject:"startState"`
			FailedState process.State `inject:"failedState"`
		}{
			StartState:  rec,
			FailedState: stubFinalState{},
		},
	)

	p, err := factory.NewFromProcessContext(process.Context{
		CurrentStateName:   "recording",
		RollbackReferences: []process.RollbackReference{{StateName: "recording"}},
	})
	require.NoError(t, err)

	// This mirrors exactly what Coordinator.Cancel does: pass the reason via Failed.
	p.Failed(context.Background(), process.CanceledByCustomerReason{PaymentCancellationReason: paymentdomain.CancellationReasonAbortedByCustomer})

	assert.True(t, rec.rolledBack)
	assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, rec.gotReason)
	assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, p.Context().CancelReason)
	assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, p.Context().FailedReason.(process.CanceledByCustomerReason).PaymentCancellationReason)
}
