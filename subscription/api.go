package subscription

import "github.com/strategicpause/slashie/actor"

type Manager interface {
	// Subscribe adds a callback when the given actor transitions to the given status.
	Subscribe(actorKey actor.Key, status actor.Status, callback Subscription)
	HandleSubscriptionsForStatus(actorKey actor.Key, status actor.Status, callback func(s Subscription))
}
