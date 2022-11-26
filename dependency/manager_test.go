package dependency

import (
	"github.com/strategicpause/slashie/actor"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	InvalidActorKey = "InvalidActorKey"
	SrcActorKey     = "SrcActorKey"
	DepActorKey     = "DepActorKey"

	DepStatus     = "DepStatus"
	SrcStatus     = "SrcStatus"
	MissingStatus = "MissingStatus"
)

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
