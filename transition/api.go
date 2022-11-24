package transition

import (
	"github.com/strategicpause/slashie/actor"
)

// Manager manages is responsible for managing transition relationships between actors.
type Manager interface {
	IsValidTransition(actorKey actor.Key, srcStatus actor.Status, depStatus actor.Status) bool
	// AddTransitionDependency is used to indicate that srcActor cannot transition srcStatus until depActor has
	// transitioned to depStatus. This will return an error if the dependency results in a invalid state (ie: a
	// circular dependency).
	AddTransitionDependency(srcActor actor.Key, srcStatus actor.Status, depActor actor.Key, depStatus actor.Status) error
	// HasTransitionDependencies returns true if the given actor has no dependencies on transitioning to the given status.
	HasTransitionDependencies(actorKey actor.Key, status actor.Status) bool
	// AddTransitionAction adds a callback function to execute to determine if the given actor can transition from
	// the srcStatus to the destStatus. If the callback fails returns an error, then the transition will not occur.
	AddTransitionAction(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status, callback Action)
	// GetTransitionActions returns a list of callbacks to execute for when the given actor wants to transition from
	// the srcStatus to the destStatus.
	GetTransitionActions(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status) []Action
	// Subscribe adds a callback when the given actor transitions to the given status.
	Subscribe(actorKey actor.Key, status actor.Status, callback Subscription)
	// GetSubscriptionsForStatus returns all subscriptions for the given actor and status.
	GetSubscriptionsForStatus(actorKey actor.Key, status actor.Status) []Subscription
	// ClearSubscriptionsForStatus
	ClearSubscriptionsForStatus(actorKey actor.Key, status actor.Status)
	// ClearDependenciesForStatus
	ClearDependenciesForStatus(actorKey actor.Key, newStatus actor.Status)
	// GetDependenciesForStatus
	GetDependenciesForStatus(actorKey actor.Key, newStatus actor.Status) []actor.Key
	StartTransition(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status, f func(a Action))
	CompleteTransitionAction(actorKey actor.Key, result error, resultFunc func(results chan error))
}
