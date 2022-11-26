package slashie

import (
	"sync"
	"testing"

	"github.com/strategicpause/slashie/logger"
	"github.com/stretchr/testify/assert"
)

func TestSingleDependencyActor(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))

	srcActor, err := NewMultiTransitionActor(1, "Src", s)
	assert.Nil(t, err)

	depActor, err := NewMultiTransitionActor(1, "Dep", s)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusA, depActor, StatusA)
	assert.Nil(t, err)

	srcActorVisited, depActorVisited := false, false
	wg := sync.WaitGroup{}
	wg.Add(2)
	err = s.Subscribe(srcActor, StatusA, func() {
		srcActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)
	err = s.Subscribe(depActor, StatusA, func() {
		depActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.UpdateStatus(srcActor, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(depActor, StatusA)
	assert.NoError(t, err)

	wg.Wait()
	assert.True(t, srcActorVisited)
	assert.True(t, depActorVisited)
}

// ActorA has 2 transitions and depends on ActorB which has 1 transition.
func TestSingleDependencyActorTwoTransitions(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))

	srcActor, err := NewMultiTransitionActor(2, "Src", s)
	assert.Nil(t, err)

	depActor, err := NewMultiTransitionActor(1, "Dep", s)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusA, depActor, StatusA)
	assert.Nil(t, err)

	srcActorStatusAVisited, srcActorStatusBVisited, depActorStatusAVisited := false, false, false
	wg := sync.WaitGroup{}
	wg.Add(3)
	err = s.Subscribe(srcActor, StatusA, func() {
		srcActorStatusAVisited = true
		wg.Done()
	})
	assert.NoError(t, err)
	err = s.Subscribe(srcActor, StatusB, func() {
		srcActorStatusBVisited = true
		wg.Done()
	})
	assert.NoError(t, err)
	err = s.Subscribe(depActor, StatusA, func() {
		depActorStatusAVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.UpdateStatus(srcActor, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(depActor, StatusA)
	assert.NoError(t, err)

	wg.Wait()
	assert.True(t, srcActorStatusAVisited)
	assert.True(t, srcActorStatusBVisited)
	assert.True(t, depActorStatusAVisited)
}

// ActorA has 2 transitions and the first transition depends on ActorB and the second transition depends on ActorC.
func TestTwoDependencyActor(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))

	srcActor, err := NewMultiTransitionActor(2, "Src", s)
	assert.Nil(t, err)

	statusADep, err := NewMultiTransitionActor(1, "DepA", s)
	assert.Nil(t, err)

	statusBDep, err := NewMultiTransitionActor(1, "DepB", s)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusA, statusADep, StatusA)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusB, statusBDep, StatusA)
	assert.Nil(t, err)

	srcActorStatusAVisited, srcActorStatusBVisited, depAActorVisited, depBActorVisited := false, false, false, false
	wg := sync.WaitGroup{}
	wg.Add(4)
	err = s.Subscribe(srcActor, StatusA, func() {
		srcActorStatusAVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(srcActor, StatusB, func() {
		srcActorStatusBVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusADep, StatusA, func() {
		depAActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusBDep, StatusA, func() {
		depBActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.UpdateStatus(srcActor, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusADep, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusBDep, StatusA)
	assert.NoError(t, err)

	wg.Wait()
	assert.True(t, srcActorStatusAVisited)
	assert.True(t, srcActorStatusBVisited)
	assert.True(t, depAActorVisited)
	assert.True(t, depBActorVisited)
}

// ActorA has 3 transitions where each transition depends on ActorB, C, & D.
func TestThreeDependencyActor(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))

	srcActor, err := NewMultiTransitionActor(3, "Src", s)
	assert.Nil(t, err)

	statusADep, err := NewMultiTransitionActor(1, "DepA", s)
	assert.Nil(t, err)

	statusBDep, err := NewMultiTransitionActor(1, "DepB", s)
	assert.Nil(t, err)

	statusCDep, err := NewMultiTransitionActor(1, "DepC", s)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusA, statusADep, StatusA)
	assert.Nil(t, err)
	err = s.AddTransitionDependency(srcActor, StatusB, statusBDep, StatusA)
	assert.Nil(t, err)
	err = s.AddTransitionDependency(srcActor, StatusC, statusCDep, StatusA)
	assert.Nil(t, err)

	srcActorStatusAVisited, srcActorStatusBVisited, srcActorStatusCVisited := false, false, false
	depAActorVisited, depBActorVisited, depCActorVisited := false, false, false
	wg := sync.WaitGroup{}
	wg.Add(6)
	// Setup subscriptions to validate we visit each expected status
	err = s.Subscribe(srcActor, StatusA, func() {
		srcActorStatusAVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(srcActor, StatusB, func() {
		srcActorStatusBVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(srcActor, StatusC, func() {
		srcActorStatusCVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusADep, StatusA, func() {
		depAActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusBDep, StatusA, func() {
		depBActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusCDep, StatusA, func() {
		depCActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.UpdateStatus(srcActor, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusADep, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusBDep, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusCDep, StatusA)
	assert.NoError(t, err)

	wg.Wait()
	assert.True(t, srcActorStatusAVisited)
	assert.True(t, srcActorStatusBVisited)
	assert.True(t, srcActorStatusCVisited)
	assert.True(t, depAActorVisited)
	assert.True(t, depBActorVisited)
	assert.True(t, depCActorVisited)
}

// ActorA has 1 transition which has 3 dependencies
func TestSingleTransitionActorThreeDependencies(t *testing.T) {
	s := NewSlashie(WithLogger(logger.NewStdOutLogger()))

	srcActor, err := NewMultiTransitionActor(1, "Src", s)
	assert.Nil(t, err)

	statusADep, err := NewMultiTransitionActor(1, "DepA", s)
	assert.Nil(t, err)

	statusBDep, err := NewMultiTransitionActor(1, "DepB", s)
	assert.Nil(t, err)

	statusCDep, err := NewMultiTransitionActor(1, "DepC", s)
	assert.Nil(t, err)

	err = s.AddTransitionDependency(srcActor, StatusA, statusADep, StatusA)
	assert.Nil(t, err)
	err = s.AddTransitionDependency(srcActor, StatusA, statusBDep, StatusA)
	assert.Nil(t, err)
	err = s.AddTransitionDependency(srcActor, StatusA, statusCDep, StatusA)
	assert.Nil(t, err)

	srcActorStatusAVisited, depAActorVisited, depBActorVisited, depCActorVisited := false, false, false, false
	wg := sync.WaitGroup{}
	wg.Add(4)
	// Setup subscriptions to validate we visit each expected status
	err = s.Subscribe(srcActor, StatusA, func() {
		srcActorStatusAVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusADep, StatusA, func() {
		depAActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusBDep, StatusA, func() {
		depBActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.Subscribe(statusCDep, StatusA, func() {
		depCActorVisited = true
		wg.Done()
	})
	assert.NoError(t, err)

	err = s.UpdateStatus(srcActor, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusADep, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusBDep, StatusA)
	assert.NoError(t, err)
	err = s.UpdateStatus(statusCDep, StatusA)
	assert.NoError(t, err)

	wg.Wait()
	assert.True(t, srcActorStatusAVisited)
	assert.True(t, depAActorVisited)
	assert.True(t, depBActorVisited)
	assert.True(t, depCActorVisited)
}
