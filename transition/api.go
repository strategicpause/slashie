package transition

import (
	"github.com/strategicpause/slashie/actor"
)

// Manager manages is responsible for managing transition relationships between actors.
type Manager interface {
	// AddTransitionAction adds a callback function to execute to determine if the given actor can transition from
	// the srcStatus to the destStatus. If the callback fails returns an error, then the transition will not occur.
	AddTransitionAction(actorKey actor.Key, srcStatus actor.Status, destStatus actor.Status, callback Action)
	// IsValidTransition returns true if the given actor is configured to transition from the srcStatus to the
	// depStatus. This is indicated by whether or not a transaction action has been added for the given source &
	// destination status.
	IsValidTransition(actorKey actor.Key, srcStatus actor.Status, depStatus actor.Status) bool
	// StartTransition will manage the actions to transition the given actor from the currentStatus to the
	// desiredStatus. Each action will be provided as a parameter to the given function.
	StartTransition(actorKey actor.Key, currentStatus actor.Status, desiredStatus actor.Status, f func(a Action))
	// CompleteTransitionAction is called when an Action has completed running. The results of all actions will be
	// sent to the given resultFunc to determine what steps to take next.
	CompleteTransitionAction(actorKey actor.Key, result error, resultFunc func(results chan error))
}
