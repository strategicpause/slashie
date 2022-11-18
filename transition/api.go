package transition

import (
	"github.com/strategicpause/slashie/actor"
)

// Manager manages is responsible for managing transition relationships between actors.
type Manager interface {
	// AddTransitionDependency is used to indicate that srcActor cannot transition srcStatus until depActor has
	// transitioned to depStatus. This will return an error if the dependency results in a invalid state (ie: a
	// circular dependency).
	AddTransitionDependency(srcActor actor.Key, srcStatus actor.Status, depActor actor.Key, depStatus actor.Status) error
	// CanTransitionToStatus returns true if the given actor has no dependencies on transitioning to the given status.
	CanTransitionToStatus(actorKey actor.Key, status actor.Status) bool
	// NotifyDependenciesOfStatus will execute the given notifFunc for each actor that is dependent on the given actor
	// when it transitions to the given newStatus.
	NotifyDependenciesOfStatus(actorKey actor.Key, newStatus actor.Status, notifFunc func(actor.Key, actor.Status))
	// AddTransitionCallback adds a callback function to execute to determine if the given actor can transition from
	// the srcStatus to the destStatus. If the callback fails returns an error, then the transition will not occur.
	AddTransitionCallback(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status, callback Callback)
	// GetTransitionCallbacks returns a list of callbacks for when the given actor transitions from the srcStatus to
	// the destStatus.
	GetTransitionCallbacks(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status) []Callback
	// Subscribe adds a callback when the given actor transitions to the given status.
	Subscribe(actorKey actor.Key, status actor.Status, callback Subscription)
	// GetSubscriptionsForStatus returns all subscriptions for the given actor and status.
	GetSubscriptionsForStatus(actorKey actor.Key, status actor.Status) []Subscription
}
