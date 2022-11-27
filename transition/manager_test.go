package transition

import (
	"fmt"
	"testing"

	"github.com/strategicpause/slashie/actor"
	"github.com/stretchr/testify/assert"
)

const (
	ActorKey        = "ActorKey"
	InvalidActorKey = "InvalidActorKey"
	SrcActorKey     = "SrcActorKey"
	DepActorKey     = "DepActorKey"

	DepStatus     = "DepStatus"
	SrcStatus     = "SrcStatus"
	DestStatus    = "DestStatus"
	MissingStatus = "MissingStatus"
)

type IsValidTransitionTest struct {
	actorKey   actor.Key
	srcStatus  actor.Status
	destStatus actor.Status
	msg        string
	result     bool
}

func TestIsValidTransition(t *testing.T) {
	mgr := NewManager()
	mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error { return nil })

	tests := []*IsValidTransitionTest{
		{actorKey: InvalidActorKey, srcStatus: SrcStatus, destStatus: DestStatus, result: false,
			msg: "return false when no transition exists for the given actor"},
		{actorKey: ActorKey, srcStatus: MissingStatus, destStatus: DestStatus, result: false,
			msg: "return false when no transition exists for the given srcStatus"},
		{actorKey: ActorKey, srcStatus: SrcStatus, destStatus: MissingStatus, result: false,
			msg: "return false when no transition exists for the given destStatus"},
		{actorKey: ActorKey, srcStatus: SrcStatus, destStatus: DestStatus, result: true,
			msg: "return true when transition exists for the given src & destStatus"},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			isValidTransition := mgr.IsValidTransition(test.actorKey, test.srcStatus, test.destStatus)
			assert.Equal(t, test.result, isValidTransition)
		})
	}
}

type StartTransitionTest struct {
	actorKey   actor.Key
	srcStatus  actor.Status
	destStatus actor.Status
	msg        string
	numActions int
}

// Actor does not exist
// Current status does not exist
// Desired status does not exist
// 1 Action is registered
// 2 Actions are registered
func TestStartTransition(t *testing.T) {
	mgr := NewManager()

	// ActorKey has one action registered
	mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error {
		return nil
	})

	// DepActorKey has two actions registered
	mgr.AddTransitionAction(DepActorKey, SrcStatus, DestStatus, func() error {
		return nil
	})
	mgr.AddTransitionAction(DepActorKey, SrcStatus, DestStatus, func() error {
		return nil
	})

	tests := []*StartTransitionTest{
		{actorKey: InvalidActorKey, srcStatus: SrcStatus, destStatus: DestStatus,
			msg: "return no actions when the actor doesn't exist", numActions: 0},
		{actorKey: ActorKey, srcStatus: MissingStatus, destStatus: DestStatus,
			msg: "return no actions when the source status doesn't exist", numActions: 0},
		{actorKey: ActorKey, srcStatus: SrcStatus, destStatus: MissingStatus,
			msg: "return no actions when the dest status doesn't exist", numActions: 0},
		{actorKey: ActorKey, srcStatus: SrcStatus, destStatus: DestStatus,
			msg: "return 1 action when one action is registered", numActions: 1},
		{actorKey: DepActorKey, srcStatus: SrcStatus, destStatus: DestStatus,
			msg: "return 2 actions when two actions are registered", numActions: 2},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			numTimesInvoked := 0
			mgr.StartTransition(test.actorKey, test.srcStatus, test.destStatus, func(a Action) {
				numTimesInvoked += 1
			})

			assert.Equal(t, test.numActions, numTimesInvoked)
		})
	}
}

type CompleteTransitionActionTest struct {
	actorKey         actor.Key
	result           error
	msg              string
	resultFuncCalled bool
}

func TestCompleteTransitionAction_OneResult(t *testing.T) {
	tests := []*CompleteTransitionActionTest{
		{actorKey: InvalidActorKey, result: nil,
			msg: "the results function should not be called if the actor does not exist", resultFuncCalled: false},
		{actorKey: ActorKey, result: nil,
			msg: "the results function should be called for a valid actor with a nil result", resultFuncCalled: true},
		{actorKey: ActorKey, result: fmt.Errorf(""),
			msg: "the results function should be called for a valid actor with a non-nil result", resultFuncCalled: true},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			mgr := NewManager()
			mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error {
				return nil
			})
			mgr.StartTransition(ActorKey, SrcStatus, DestStatus, func(a Action) {})

			resultFuncCalled := false
			mgr.CompleteTransitionAction(test.actorKey, test.result, func(results chan error) {
				resultFuncCalled = true

				assert.Equal(t, test.result, <-results)
			})

			assert.Equal(t, test.resultFuncCalled, resultFuncCalled)
		})
	}
}

func TestCompleteTransitionAction_TwoResults(t *testing.T) {
	mgr := NewManager()
	// Register two transition actions
	mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error {
		return nil
	})
	mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error {
		return nil
	})
	mgr.StartTransition(ActorKey, SrcStatus, DestStatus, func(a Action) {})
	// Complete the first action. The function should not be called yet.
	resultFuncCalled := false
	mgr.CompleteTransitionAction(ActorKey, nil, func(results chan error) {
		resultFuncCalled = true
	})
	assert.False(t, resultFuncCalled)
	// Complet the section action. Verify the given function is called.
	mgr.CompleteTransitionAction(ActorKey, nil, func(results chan error) {
		resultFuncCalled = true

		assert.Nil(t, <-results)
	})
	assert.True(t, resultFuncCalled)
}
