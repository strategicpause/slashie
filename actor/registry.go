package actor

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
	actorKey := actor.GetKey()
	a.actors[actorKey] = actor

	return actorKey
}

func (r *registry) GetActor(actorKey Key) (Actor, bool) {
	a, ok := r.actors[actorKey]

	return a, ok
}

func (r *registry) IsRegistered(actor Actor) bool {
	_, ok := r.actors[actor.GetKey()]

	return ok
}
