package dependency

import (
	"github.com/strategicpause/slashie/actor"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

const (
	InvalidActorKey = "InvalidActorKey"
	SrcActorKey     = "SrcActorKey"
	DepActorKey     = "DepActorKey"
	OtherActorKey   = "OtherActorKey"

	ActorA = "ActorA"
	ActorB = "ActorB"
	ActorC = "ActorC"
	ActorD = "ActorD"

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

func TestNotifyDependenciesOfStatus_OneDependency(t *testing.T) {
	mgr := NewManager()
	err := mgr.AddTransitionDependency(SrcActorKey, SrcStatus, DepActorKey, DepStatus)
	assert.Nil(t, err)

	isSrcActor := false
	timesCalled := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	mgr.NotifyDependenciesOfStatus(DepActorKey, DepStatus, func(actorKey actor.Key) {
		isSrcActor = actorKey == SrcActorKey
		timesCalled += 1
		wg.Done()
	})
	wg.Wait()

	assert.True(t, isSrcActor)
	assert.Equal(t, 1, timesCalled)
}

func TestNotifyDependenciesOfStatus_TwoDependency(t *testing.T) {
	mgr := NewManager()
	err := mgr.AddTransitionDependency(SrcActorKey, SrcStatus, DepActorKey, DepStatus)
	assert.Nil(t, err)
	err = mgr.AddTransitionDependency(SrcActorKey, SrcStatus, OtherActorKey, DepStatus)
	assert.Nil(t, err)

	isSrcActor := false
	timesCalled := 0
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Since there are two dependencies, the srcActor should only be notified once both dependencies have
	// been removed.
	mgr.NotifyDependenciesOfStatus(DepActorKey, DepStatus, func(actorKey actor.Key) {
		isSrcActor = actorKey == SrcActorKey
		timesCalled += 1
		wg.Done()
	})
	mgr.NotifyDependenciesOfStatus(OtherActorKey, DepStatus, func(actorKey actor.Key) {
		isSrcActor = actorKey == SrcActorKey
		timesCalled += 1
		wg.Done()
	})
	wg.Wait()

	assert.True(t, isSrcActor)
	assert.Equal(t, 1, timesCalled)
}

// Verify that the code can catch circular dependencies between the same actor:
// ActorA -> ActorA
func TestAddTransitionDependency_SameActorCircularDependency(t *testing.T) {
	mgr := NewManager()
	err := mgr.AddTransitionDependency(ActorA, SrcStatus, ActorA, SrcStatus)

	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between two actors:
// ActorA -> ActorB -> ActorA
func TestAddTransitionDependency_SimpleCircularDependency(t *testing.T) {
	mgr := NewManager()

	err := mgr.AddTransitionDependency(ActorA, SrcStatus, ActorB, DepStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorB, DepStatus, ActorA, SrcStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between 3+ actors:
// ActorA -> ActorB -> ActorC -> ActorA
func TestAddTransitionDependency_CircularDependency(t *testing.T) {
	mgr := NewManager()

	err := mgr.AddTransitionDependency(ActorA, SrcStatus, ActorB, DepStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorB, DepStatus, ActorC, DepStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorC, DepStatus, ActorA, SrcStatus)
	assert.Error(t, err)
}

// Verify that the code does result in a false positive when two actors have a dependency
// on the same actor. For example:
//
//	/ B \
//
// A     D
//
//	\ C /
func TestAddTransitionDependency_NonCircularDependency(t *testing.T) {
	mgr := NewManager()

	err := mgr.AddTransitionDependency(ActorA, SrcStatus, ActorB, SrcStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorA, SrcStatus, ActorC, SrcStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorB, SrcStatus, ActorD, SrcStatus)
	assert.NoError(t, err)

	err = mgr.AddTransitionDependency(ActorC, SrcStatus, ActorD, SrcStatus)
	assert.NoError(t, err)
}
