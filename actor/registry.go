package actor

import (
	"fmt"
)

type registry struct {
	// actors are used to resolve an ActorKey to an Actor
	actors map[Key]Actor
}

func NewRegistry() Registry {
	return &registry{
		actors: map[Key]Actor{},
	}
}

func (a *registry) RegisterActor(actor Actor) Key {
	actorKey := a.getActorKey(actor)

	a.actors[actorKey] = actor

	return actorKey
}

// getActorKey creates a globally unique identifier for an actor. While the ActorId may not be globally unique,
// the unique ActorType acts a namespace such that the returned key is globally unique.
func (a *registry) getActorKey(actor Actor) Key {
	actorKey := fmt.Sprintf("%s-%s", actor.GetType(), actor.GetId())
	return Key(actorKey)
}

func (a *registry) GetActorKey(actor Actor) (Key, bool) {
	actorKey := a.getActorKey(actor)
	_, ok := a.actors[actorKey]
	if !ok {
		return "", ok
	}

	return actorKey, ok
}

func (a *registry) GetActor(actorKey Key) Actor {
	return a.actors[actorKey]
}
