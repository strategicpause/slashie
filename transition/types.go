package transition

import (
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

// ActionsByStatus
type ActionsByStatus map[actor.Status]map[actor.Status][]Action
