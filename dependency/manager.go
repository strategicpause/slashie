package dependency

import (
	"errors"
	"github.com/strategicpause/slashie/actor"
)

type manager struct {
	// transitionDependenciesByActor has a structure of map[actor.Key][DesiredStatus][DependentActor][Status] and
	// A given actor cannot transition to the DesiredStatus until all DependentActors have transitioned
	// to their given status. When all the transition dependencies for a given actor by state have been satisfied,
	//then the callback will be called and removed from the map.
	transitionDependenciesByActor map[actor.Key]Dependencies
	// reverseDependencies
	reverseDependencies map[actor.Key]map[actor.Status]map[actor.Key]actor.Status
}

func NewManager() Manager {
	return &manager{
		transitionDependenciesByActor: map[actor.Key]Dependencies{},
		reverseDependencies:           map[actor.Key]map[actor.Status]map[actor.Key]actor.Status{},
	}
}

func (t *manager) HasTransitionDependencies(actorKey actor.Key, status actor.Status) bool {
	if _, ok := t.transitionDependenciesByActor[actorKey]; !ok {
		return false
	}
	transitionDependencies := t.transitionDependenciesByActor[actorKey]
	if _, ok := transitionDependencies[status]; !ok {
		return false
	}
	// There are still dependencies
	if len(transitionDependencies[status]) > 0 {
		return true
	}

	return false
}

func (t *manager) NotifyDependenciesOfStatus(actorKey actor.Key, newStatus actor.Status, callback func(actor.Key)) {
	deps := t.reverseDependencies[actorKey][newStatus]
	for depActorKey, depActorStatus := range deps {
		delete(t.transitionDependenciesByActor[depActorKey][depActorStatus], actorKey)
		// This means the given dependent actor no longer has any dependencies and can transition to the desired status.
		if len(t.transitionDependenciesByActor[depActorKey][depActorStatus]) == 0 {
			callback(depActorKey)
		}
	}
	delete(t.reverseDependencies[actorKey], newStatus)
}

func (t *manager) AddTransitionDependency(srcActor actor.Key, srcStatus actor.Status, depActor actor.Key, depStatus actor.Status) error {
	if _, ok := t.transitionDependenciesByActor[srcActor]; !ok {
		t.transitionDependenciesByActor[srcActor] = Dependencies{}
	}
	transitionDependencies := t.transitionDependenciesByActor[srcActor]
	if _, ok := transitionDependencies[srcStatus]; !ok {
		transitionDependencies[srcStatus] = map[actor.Key]actor.Status{}
	}
	dependencies := transitionDependencies[srcStatus]
	dependencies[depActor] = depStatus
	err := t.validateTransitionDependencies(srcActor, srcStatus)
	// If this results in an invalid state, then undo the previous action
	if err != nil {
		transitionDependencies[srcStatus] = map[actor.Key]actor.Status{}
	} else {
		if _, ok := t.reverseDependencies[depActor]; !ok {
			t.reverseDependencies[depActor] = make(map[actor.Status]map[actor.Key]actor.Status)
		}
		if _, ok := t.reverseDependencies[depActor][depStatus]; !ok {
			t.reverseDependencies[depActor][depStatus] = make(map[actor.Key]actor.Status)
		}
		t.reverseDependencies[depActor][depStatus][srcActor] = srcStatus
	}

	return err
}

func (t *manager) validateTransitionDependencies(srcActor actor.Key, srcStatus actor.Status) error {
	visited := map[string]bool{}
	var queue []*ActorStatusKey
	queue = append(queue, &ActorStatusKey{
		actorKey: srcActor,
		status:   srcStatus,
	})
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		currKey := curr.String()

		visited[currKey] = true
		// Add dependencies to the queue
		dependencies := t.transitionDependenciesByActor[curr.actorKey][curr.status]
		for depActor, status := range dependencies {
			key := &ActorStatusKey{
				actorKey: depActor,
				status:   status,
			}
			if _, ok := visited[key.String()]; ok {
				return errors.New("already visited")
			}
			queue = append(queue, key)
		}
	}
	return nil
}
