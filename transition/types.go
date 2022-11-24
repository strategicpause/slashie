package transition

import (
	"fmt"
	"github.com/strategicpause/slashie/actor"
)

// TransitionAction encapsulates a callback for when an actor transitions from SrcStatus to DestStatus.
type TransitionAction struct {
	SrcStatus  actor.Status
	DestStatus actor.Status
	Action     Action
}

// Action actor to register a callback to execute when i
type Action func() error

// Subscription is a function to execute after an actor transitions to a given state.
type Subscription func()

// Dependencies
type Dependencies map[actor.Status]map[actor.Key]actor.Status

// ActionsByStatus
type ActionsByStatus map[actor.Status]map[actor.Status][]Action

// SubscriptionsByStatus
type SubscriptionsByStatus map[actor.Status][]Subscription

// ActorStatusKey
type ActorStatusKey struct {
	actorKey actor.Key
	status   actor.Status
}

func (a *ActorStatusKey) String() string {
	return fmt.Sprintf("%s-%s", a.actorKey, a.status)
}
