package slashie

import (
	"errors"
	"github.com/strategicpause/slashie/actor"
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
	tm := NewSlashie()
	srcActor := NewBasicActor("SourceActor", "ActorA", tm)
	tm.AddActor(srcActor, NoneStatus, StoppedStatus)
	depActor := NewBasicActor("DependentActor", "ActorB", tm)
	tm.AddActor(depActor, NoneStatus, StoppedStatus)

	// ActorA cannot transition to READY until ActorB first transitions to READY.
	err := tm.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
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

func NewBasicActor(actorType actor.Type, actorId actor.Id, tm Slashie) actor.Actor {
	basicActor := &actor.BasicActor{
		Type:    actorType,
		Id:      actorId,
		Mailbox: make(chan func(), 10),
	}
	go basicActor.Init()

	tm.AddActor(basicActor, NoneStatus, StoppedStatus)

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
	err := tm.AddTransitionCallback(basicActor, NoneStatus, ReadyStatus, callback)
	assert.NoError(t, err)

	err = tm.UpdateStatus(basicActor, ReadyStatus)
	assert.NoError(t, err)

	// Indicates the callback was properly executed
	assert.True(t, <-callbackChan)
}

func TestUpdateStatus_NoTransitionCallback(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	wg := sync.WaitGroup{}
	// Add a subscription, so  that we can wait
	err := tm.Subscribe(basicActor, ReadyStatus, func() {
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
	basicActor := &actor.BasicActor{
		Type: "actor",
		Id:   "id",
	}

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

	err := tm.AddTransitionCallback(basicActor, NoneStatus, ReadyStatus, func() error {
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
	tm := NewSlashie()

	srcActor := &actor.BasicActor{
		Type: "Actor",
		Id:   "src",
	}
	depActor := NewBasicActor("Actor", "dep", tm)

	err := tm.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
	assert.Error(t, err)
}

func TestAddTransitionDependency_UnknownDependentActor(t *testing.T) {
	tm := NewSlashie()

	srcActor := NewBasicActor("Actor", "src", tm)
	depActor := &actor.BasicActor{
		Type: "Actor",
		Id:   "dep",
	}

	err := tm.AddTransitionDependency(srcActor, ReadyStatus, depActor, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between the same actor
func TestAddTransitionDependency_SameActorCircularDependency(t *testing.T) {
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorA, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between two actorMap
func TestAddTransitionDependency_SimpleCircularDependency(t *testing.T) {
	tm := NewSlashie()
	actorA := NewBasicActor("Actor", "ActorA", tm)
	actorB := NewBasicActor("Actor", "ActorB", tm)

	err := tm.AddTransitionDependency(actorA, ReadyStatus, actorB, ReadyStatus)
	assert.NoError(t, err)

	err = tm.AddTransitionDependency(actorB, ReadyStatus, actorA, ReadyStatus)
	assert.Error(t, err)
}

// Verify that the code can catch circular dependencies between 3+ actorMap
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

// Verify that the code does result in a false positive when two actorMap have a dependency
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

	err := tm.AddTransitionCallback(basicActor, ReadyStatus, ReadyStatus, func() error { return nil })
	assert.Error(t, err)
}

// Verify that a callback cannot be added to the initial status.
func TestAddTransitionCallback_StartState(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionCallback(basicActor, ReadyStatus, NoneStatus, func() error { return nil })
	assert.Error(t, err)
}

// Verify that a callback cannot be added from the terminal status.
func TestAddTransitionCallback_TermState(t *testing.T) {
	tm := NewSlashie()
	basicActor := NewBasicActor("Actor", "ActorA", tm)

	err := tm.AddTransitionCallback(basicActor, StoppedStatus, NoneStatus, func() error { return nil })
	assert.Error(t, err)
}
