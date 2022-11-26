package dependency

import (
	"fmt"
	"github.com/strategicpause/slashie/actor"
)

// Dependencies
type Dependencies map[actor.Status]map[actor.Key]actor.Status

// ActorStatusKey
type ActorStatusKey struct {
	actorKey actor.Key
	status   actor.Status
}

func (a *ActorStatusKey) String() string {
	return fmt.Sprintf("%s-%s", a.actorKey, a.status)
}
