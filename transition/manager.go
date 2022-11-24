package transition

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
	// transitionActionsByActor
	transitionActionsByActor map[actor.Key]ActionsByStatus
	// subscriptionsForActor
	subscriptionsForActor  map[actor.Key]SubscriptionsByStatus
	transitionsByActorChan map[actor.Key]chan error
}

func NewManager() Manager {
	return &manager{
		transitionDependenciesByActor: map[actor.Key]Dependencies{},
		reverseDependencies:           map[actor.Key]map[actor.Status]map[actor.Key]actor.Status{},
		transitionActionsByActor:      map[actor.Key]ActionsByStatus{},
		subscriptionsForActor:         map[actor.Key]SubscriptionsByStatus{},
		transitionsByActorChan:        map[actor.Key]chan error{},
	}
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

func (t *manager) GetTransitionActions(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status) []Action {
	if _, ok := t.transitionActionsByActor[actorKey]; !ok {
		return []Action{}
	}
	transitionCallbacks := t.transitionActionsByActor[actorKey]
	if _, ok := transitionCallbacks[currentStatus]; !ok {
		return []Action{}
	}

	return t.transitionActionsByActor[actorKey][currentStatus][desiredStatus]
}

func (t *manager) GetSubscriptionsForStatus(actorKey actor.Key, status actor.Status) []Subscription {
	if _, ok := t.subscriptionsForActor[actorKey]; !ok {
		return []Subscription{}
	}
	return t.subscriptionsForActor[actorKey][status]
}

func (t *manager) ClearSubscriptionsForStatus(actorKey actor.Key, status actor.Status) {
	if _, ok := t.subscriptionsForActor[actorKey]; !ok {
		return
	}
	delete(t.subscriptionsForActor[actorKey], status)
}

func (t *manager) GetDependenciesForStatus(actorKey actor.Key, newStatus actor.Status) []actor.Key {
	var actorsToNotify []actor.Key
	deps := t.reverseDependencies[actorKey][newStatus]
	for depActorKey := range deps {
		actorsToNotify = append(actorsToNotify, depActorKey)
	}

	return actorsToNotify
}

func (t *manager) ClearDependenciesForStatus(actorKey actor.Key, newStatus actor.Status) {
	deps := t.reverseDependencies[actorKey][newStatus]
	for depActorKey, depActorStatus := range deps {
		delete(t.transitionDependenciesByActor[depActorKey][depActorStatus], actorKey)
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

func (t *manager) Subscribe(actorKey actor.Key, status actor.Status, callback Subscription) {
	if _, ok := t.subscriptionsForActor[actorKey]; !ok {
		t.subscriptionsForActor[actorKey] = SubscriptionsByStatus{}
	}
	subscriptionsByStatus := t.subscriptionsForActor[actorKey]

	subscriptionsByStatus[status] = append(subscriptionsByStatus[status], callback)
}

func (t *manager) StartTransition(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status, f func(a Action)) {
	actions := t.GetTransitionActions(actorKey, currentStatus, desiredStatus)
	numActions := len(actions)
	t.transitionsByActorChan[actorKey] = make(chan error, numActions)

	for _, action := range actions {
		f(action)
	}
}

func (t *manager) CompleteTransitionAction(actorKey actor.Key, result error, resultFunc func(results chan error)) {
	results := t.transitionsByActorChan[actorKey]
	results <- result
	// If the channel has all results, then execute resultFunc with the results.
	if len(results) == cap(results) {
		close(results)
		delete(t.transitionsByActorChan, actorKey)

		resultFunc(results)
	}
}
