package transition

import (
	"github.com/strategicpause/slashie/actor"
)

type manager struct {

	// transitionActionsByActor
	transitionActionsByActor map[actor.Key]ActionsByStatus
	transitionsByActorChan   map[actor.Key]chan error
}

func NewManager() Manager {
	return &manager{
		transitionActionsByActor: map[actor.Key]ActionsByStatus{},
		transitionsByActorChan:   map[actor.Key]chan error{},
	}
}

func (t *manager) AddTransitionAction(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status, callback Action) {
	if _, ok := t.transitionActionsByActor[actorKey]; !ok {
		t.transitionActionsByActor[actorKey] = ActionsByStatus{}
	}
	transitionCallbacks := t.transitionActionsByActor[actorKey]
	if _, ok := transitionCallbacks[srcStatus]; !ok {
		transitionCallbacks[srcStatus] = map[actor.Status][]Action{}
	}
	transitionCallbacks[srcStatus][destStatus] = append(transitionCallbacks[srcStatus][destStatus], callback)
}

func (t *manager) IsValidTransition(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status) bool {
	if _, ok := t.transitionActionsByActor[actorKey]; !ok {
		return false
	}
	transitionCallbacks := t.transitionActionsByActor[actorKey]
	if _, ok := transitionCallbacks[srcStatus]; !ok {
		return false
	}

	_, ok := t.transitionActionsByActor[actorKey][srcStatus][destStatus]
	return ok
}

func (t *manager) StartTransition(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status, f func(a Action)) {
	if _, ok := t.transitionActionsByActor[actorKey]; !ok {
		return
	}

	transitionCallbacks := t.transitionActionsByActor[actorKey]
	if _, ok := transitionCallbacks[currentStatus]; !ok {
		return
	}

	actions := transitionCallbacks[currentStatus][desiredStatus]
	numActions := len(actions)
	t.transitionsByActorChan[actorKey] = make(chan error, numActions)

	for _, action := range actions {
		f(action)
	}
}

func (t *manager) CompleteTransitionAction(actorKey actor.Key, result error, resultFunc func(results chan error)) {
	results, ok := t.transitionsByActorChan[actorKey]
	if !ok {
		return
	}
	results <- result
	// If the channel has all results, then execute resultFunc with the results.
	if len(results) == cap(results) {
		close(results)
		delete(t.transitionsByActorChan, actorKey)

		resultFunc(results)
	}
}
