package transition

import (
	"github.com/strategicpause/slashie/actor"
	"github.com/stretchr/testify/assert"
	"testing"
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

type HasTransitionDependenciesTests struct {
	actorKey actor.Key
	status   actor.Status
	msg      string
	result   bool
}

func TestHasTransitionDependencies(t *testing.T) {
	mgr := NewManager()
	err := mgr.AddTransitionDependency(SrcActorKey, SrcStatus, DepActorKey, DepStatus)
	assert.Nil(t, err)

	tests := []*HasTransitionDependenciesTests{
		{actorKey: InvalidActorKey, status: SrcStatus, result: false,
			msg: "return false when no dependencies exist for the given actor"},
		{actorKey: SrcActorKey, status: MissingStatus, result: false,
			msg: "return false when no dependencies exist for the given status"},
		{actorKey: SrcActorKey, status: SrcStatus, result: true,
			msg: "return true when there are dependencies for the given actor & status"},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			hasDependencies := mgr.HasTransitionDependencies(test.actorKey, test.status)
			assert.Equal(t, test.result, hasDependencies)
		})
	}
}

//type GetTransitionActionsTest struct {
//	actorKey       actor.Key
//	srcStatus      actor.Status
//	destStatus     actor.Status
//	msg            string
//	numTransitions int
//}
//
//func TestGetTransitionActions(t *testing.T) {
//	mgr := NewManager()
//	mgr.AddTransitionAction(ActorKey, SrcStatus, DestStatus, func() error { return nil })
//
//	tests := []*GetTransitionActionsTest{
//		{actorKey: },
//	}
//}
