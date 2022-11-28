package slashie

import (
	"errors"
	"github.com/strategicpause/slashie/actor"
	"github.com/strategicpause/slashie/transition"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

const (
	NoneStatus    actor.Status = "NONE"
	ReadyStatus   actor.Status = "READY"
	StoppedStatus actor.Status = "STOPPED"
)

func TestAddTransitionDependency(t *testing.T) {
	// Setup actors
	tm := NewSlashie(WithMailboxSize(DefaultMailboxSize))
	srcActor := NewBasicActor("SourceActor", "ActorA", tm)
	tm.AddActor(srcActor, NoneStatus, StoppedStatus)
	depActor := NewBasicActor("DependentActor", "ActorB", tm)
	tm.AddActor(depActor, NoneStatus, StoppedStatus)

	// Setup basic dependencies from NONE -> READY for each actor.
	err := tm.AddTransitionAction(srcActor, NoneStatus, ReadyStatus, func() error { return nil })
	assert.NoError(t, err)
	err = tm.AddTransitionAction(depActor, NoneStatus, ReadyStatus, func() error { return nil })
	assert.NoError(t, err)

	// ActorA cannot transition to READY until ActorB first transitions to READY.
	err = tm.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
	assert.NoError(t, err)
	// Update the status of ActorA to be READY
	err = tm.UpdateStatus(srcActor, ReadyStatus)
	assert.NoError(t, err)

	// This will be called when the srcActor successfully transitions to Ready
	ch := make(chan bool)
	cb := func() {
		ch <- true
	}
	err = tm.Subscribe(srcActor, ReadyStatus, cb)
	assert.NoError(t, err)

	// Verify that status of ActorA is still NONE since depActor still has not transitioned to Ready.
	srcStatus := tm.GetStatus(srcActor)
	assert.Equal(t, NoneStatus, srcStatus)

	// Update dependent Actor
	err = tm.UpdateStatus(depActor, ReadyStatus)
	assert.NoError(t, err)

	// Wait until the callback has been called.
	<-ch

	// Verify that the status of ActorA has also transitioned
	srcStatus = tm.GetStatus(srcActor)
	assert.Equal(t, ReadyStatus, srcStatus)
}

func NewBasicActor(actorType actor.Type, actorId actor.Id, s Slashie) actor.Actor {
	basicActor := actor.NewBasicActor(actorType, actorId)

	s.AddActor(basicActor, NoneStatus, StoppedStatus)

	return basicActor
}

func TestAddTransitionCallback(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	callbackChan := make(chan bool)
	callback := func() error {
		callbackChan <- true
		close(callbackChan)

		return nil
	}
	err := tm.AddTransitionAction(basicActor, NoneStatus, ReadyStatus, callback)
	assert.NoError(t, err)

	err = tm.UpdateStatus(basicActor, ReadyStatus)
	assert.NoError(t, err)

	// Indicates the callback was properly executed
	assert.True(t, <-callbackChan)
}

func TestUpdateStatus_NoTransitionCallback(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionAction(basicActor, NoneStatus, ReadyStatus, func() error { return nil })
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	// Add a subscription, so  that we can wait
	err = tm.Subscribe(basicActor, ReadyStatus, func() {
		wg.Done()
	})
	assert.NoError(t, err)

	wg.Add(1)
	err = tm.UpdateStatus(basicActor, ReadyStatus)
	assert.NoError(t, err)
	wg.Wait()

	status := tm.GetStatus(basicActor)
	assert.Equal(t, ReadyStatus, status)
}

func TestUpdateStatus_UnknownActor(t *testing.T) {
	tm := NewSlashie()
	basicActor := actor.NewBasicActor("actor", "id")

	err := tm.UpdateStatus(basicActor, ReadyStatus)

	assert.Error(t, err)
}

func TestGetStatus(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	status := tm.GetStatus(basicActor)
	assert.Equal(t, NoneStatus, status)
}

func TestAddActor(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	// Get default Status
	status := tm.GetStatus(basicActor)

	assert.Equal(t, NoneStatus, status)
}

func TestTransitionError(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionAction(basicActor, NoneStatus, ReadyStatus, func() error {
		return errors.New("failed to transition")
	})
	assert.NoError(t, err)

	ch := make(chan bool)
	err = tm.Subscribe(basicActor, StoppedStatus, func() {
		ch <- true
		close(ch)
	})
	assert.NoError(t, err)

	err = tm.UpdateStatus(basicActor, ReadyStatus)
	assert.NoError(t, err)
	<-ch

	status := tm.GetStatus(basicActor)
	assert.Equal(t, StoppedStatus, status)
}

func TestAddTransitionDependency_UnknownSourceActor(t *testing.T) {
	s := NewSlashie()

	srcActor := actor.NewBasicActor("Actor", "src")
	depActor := NewBasicActor("Actor", "dep", s)

	err := s.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
	assert.Error(t, err)
}

func TestAddTransitionDependency_UnknownDependentActor(t *testing.T) {
	tm := NewSlashie()

	srcActor := NewBasicActor("Actor", "src", tm)
	depActor := actor.NewBasicActor("Actor", "dep")

	err := tm.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between the same actor:
// ActorA -> ActorA
func TestAddTransitionDependency_SameActorCircularDependency(t *testing.T) {
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorA, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between two actors:
// ActorA -> ActorB -> ActorA
func TestAddTransitionDependency_SimpleCircularDependency(t *testing.T) {
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)
	actorB := NewBasicActor("Actor", "ActorB", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorB, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorB, ReadyStatus, actorA, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between 3+ actors:
// ActorA -> ActorB -> ActorC -> ActorA
func TestAddTransitionDependency_CircularDependency(t *testing.T) {
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)
	actorB := NewBasicActor("Actor", "ActorB", tm)
	actorC := NewBasicActor("Actor", "ActorC", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorB, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorB, ReadyStatus, actorC, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorC, ReadyStatus, actorA, ReadyStatus)
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
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)
	actorB := NewBasicActor("Actor", "ActorB", tm)
	actorC := NewBasicActor("Actor", "ActorC", tm)
	actorD := NewBasicActor("Actor", "ActorD", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorB, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorA, ReadyStatus, actorC, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorB, ReadyStatus, actorD, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorC, ReadyStatus, actorD, ReadyStatus)
	assert.NoError(t, err)
}

// Verify that a callback cannot be added for a transition to the same status.
func TestAddTransitionCallback_SameState(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionAction(basicActor, ReadyStatus, ReadyStatus, func() error { return nil })
	assert.Error(t, err)
}

// Verify that a callback cannot be added to the initial status.
func TestAddTransitionCallback_StartState(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionAction(basicActor, ReadyStatus, NoneStatus, func() error { return nil })
	assert.Error(t, err)
}

// Verify that a callback cannot be added from the terminal status.
func TestAddTransitionCallback_TermState(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionAction(basicActor, StoppedStatus, NoneStatus, func() error { return nil })
	assert.Error(t, err)
}

func TestAddTransitionActions(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionActions(basicActor, []*transition.TransitionAction{
		{SrcStatus: NoneStatus, DestStatus: ReadyStatus, Action: func() error {
			return nil
		}},
		{SrcStatus: ReadyStatus, DestStatus: StoppedStatus, Action: func() error {
			return nil
		}},
	})

	assert.Nil(t, err)
}

func TestAddTransitionActions_IllegalTransition(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionActions(basicActor, []*transition.TransitionAction{
		{SrcStatus: NoneStatus, DestStatus: ReadyStatus, Action: func() error {
			return nil
		}},
		{SrcStatus: ReadyStatus, DestStatus: NoneStatus, Action: func() error {
			return nil
		}},
	})

	assert.Error(t, err)
}
