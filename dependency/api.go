package dependency

import "github.com/strategicpause/slashie/actor"

type Manager interface {
	// AddTransitionDependency is used to indicate that srcActor cannot transition srcStatus until depActor has
	// transitioned to depStatus. This will return an error if the dependency results in a invalid state (ie: a
	// circular dependency).
	AddTransitionDependency(srcActor actor.Key, srcStatus actor.Status, depActor actor.Key, depStatus actor.Status) error
	// HasTransitionDependencies returns true if the given actor has no dependencies on transitioning to the given status.
	HasTransitionDependencies(actorKey actor.Key, status actor.Status) bool
	NotifyDependenciesOfStatus(actorKey actor.Key, newStatus actor.Status, callback func(actor.Key))
}
