package slashie

import (
	"github.com/strategicpause/slashie/actor"
	"github.com/strategicpause/slashie/transition"
)

// Slashie manages all callbacks and dependencies which establish relationships between the different actors.
type Slashie interface {
	// AddActor will register an Actor with the Slashie. This will also specify both the initial status
	// and terminal status for the given actor.
	AddActor(actor actor.Actor, initStatus actor.Status, terminalStatus actor.Status)
	// AddTransitionDependency will add a dependency on the srcActor transitioning to srcStatus until destActor
	// transitions to destStatus.
	AddTransitionDependency(srcActor actor.Actor, srcStatus actor.Status, depActor actor.Actor, depStatus actor.Status) error
	// AddTransitionAction will register a callback function which will be called before the given actor
	// transitions from srcStatus to destStatus.
	AddTransitionAction(actor actor.Actor, srcStatus actor.Status, destStatus actor.Status, callback transition.Action) error
	// AddTransitionActions registers multiple transition callbacks for a given Actor.
	AddTransitionActions(actor actor.Actor, transitionCallbacks []*transition.TransitionAction) error
	// UpdateStatus indicates that the given actor wants to transition to the desiredStatus. Once all of an actor's
	// dependencies have reached their desired state, then transition callbacks will be called for that actor. Upon
	// successful completion of transition callbacks, then the actor will successfully move to the desired status.
	// Once transitioning has completed, then all subscription callbacks will be called.
	UpdateStatus(actor actor.Actor, desiredStatus actor.Status) error
	// GetStatus returns the current known status for an Actor.
	GetStatus(actor actor.Actor) actor.Status
	// Subscribe allows anyone to register a callback function to execute once the given actor has transitioned
	// to the given status.
	Subscribe(actor actor.Actor, status actor.Status, callback transition.Subscription) error
}
