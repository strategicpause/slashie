package transition

import (
	"github.com/strategicpause/slashie/actor"
)

// Manager manages is responsible for managing transition relationships between actors.
type Manager interface {
	IsValidTransition(actorKey actor.Key, srcStatus actor.Status, depStatus actor.Status) bool
	// AddTransitionAction adds a callback function to execute to determine if the given actor can transition from
	// the srcStatus to the destStatus. If the callback fails returns an error, then the transition will not occur.
	AddTransitionAction(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status, callback Action)
	StartTransition(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status, f func(a Action))
	CompleteTransitionAction(actorKey actor.Key, result error, resultFunc func(results chan error))
}
