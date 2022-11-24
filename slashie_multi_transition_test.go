package slashie

import (
	"sync"
	"testing"

	"github.com/strategicpause/slashie/actor"
	"github.com/strategicpause/slashie/logger"
	"github.com/stretchr/testify/assert"
)

const (
	StatusInit = "Init"
	StatusA    = "A"
	StatusB    = "B"
	StatusC    = "C"

	ActorType      = "TestActor"
	NumTransitions = 3
)

var (
	Statuses = []actor.Status{
		StatusA,
		StatusB,
		StatusC,
	}
)

type MultiTransitionActor struct {
	*actor.BasicActor
	s Slashie
}

// NewMultiTransitionActor creates an actor with n transitions. Each transition will transition to the next until
// it reaches the terminal status.
func NewMultiTransitionActor(n int, id actor.Id, s Slashie) (actor.Actor, error) {
	a := MultiTransitionActor{
		BasicActor: actor.NewBasicActor(ActorType, id),
		s:          s,
	}
	terminalStatus := Statuses[n-1]
	s.AddActor(a, StatusInit, terminalStatus)

	prevStatus := actor.Status(StatusInit)
	for i := 0; i < n; i++ {
		currStatus := Statuses[i]
		var cb func() error
		if i < len(Statuses)-1 {
			nextStatus := Statuses[i+1]
			cb = func() error {
				return s.UpdateStatus(a, nextStatus)
			}
		} else {
			cb = func() error {
				return nil
			}
		}
		err := s.AddTransitionAction(a, prevStatus, currStatus, cb)
		if err != nil {
			return nil, err
		}
		prevStatus = currStatus
	}

	return a, nil
}

// TestMultiTransitionActor will validate an actor can transition from an initial state to
// the terminal state.
func TestMultiTransitionActor(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))
	a, err := NewMultiTransitionActor(NumTransitions, "Id", s)
	assert.Nil(t, err)

	// This will be used to track that we have actually visited all three states.
	visitedA, visitedB, visitedC := false, false, false
	wg := sync.WaitGroup{}
	wg.Add(NumTransitions)

	// Subscribe to each state to verify it has been visited.
	err = s.Subscribe(a, StatusA, func() {
		visitedA = true
		wg.Done()
	})
	assert.Nil(t, err)

	err = s.Subscribe(a, StatusB, func() {
		visitedB = true
		wg.Done()
	})
	assert.Nil(t, err)

	err = s.Subscribe(a, StatusC, func() {
		visitedC = true
		wg.Done()
	})
	assert.Nil(t, err)

	err = s.UpdateStatus(a, StatusA)
	assert.Nil(t, err)

	// Block until we have visited all states
	wg.Wait()

	assert.True(t, visitedA)
	assert.True(t, visitedB)
	assert.True(t, visitedC)
}
