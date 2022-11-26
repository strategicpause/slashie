package dependency

import (
	"fmt"
	"github.com/strategicpause/slashie/actor"
)

// ActorStatusKey is used to determine if cycles exist in the dependency map.
type ActorStatusKey struct {
	actorKey actor.Key
	status   actor.Status
}

func (a *ActorStatusKey) String() string {
	return fmt.Sprintf("%s-%s", a.actorKey, a.status)
}
