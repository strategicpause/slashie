package transition

import (
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
