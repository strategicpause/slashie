package actor

// Actor
type Actor interface {
	// GetType returns the type of Actor
	GetType() Type
	// GetId
	GetId() Id
	// Notify
	Notify(callback func())
}

type Registry interface {
	// RegisterActor will register an Actor and return the corresponding ActorKey.
	RegisterActor(actor Actor) Key
	// GetActorKey will return the key for the given Actor if it has been registered.
	GetActorKey(actor Actor) (Key, bool)
	// GetActor will return an Actor, if it exists, given an ActorKey.
	GetActor(actorKey Key) Actor
}

type StatusManager interface {
	//
	InitializeActor(actorKey Key, initStatus Status, terminalStatus Status)
	IsValidSubscriptionStatus(actorKey Key, status Status) bool
	IsValidTransitionStatus(actorKey Key, srcStatus Status, destStatus Status) bool
	GetKnownStatus(actorKey Key) Status
	SetKnownStatus(actorKey Key, knownStatus Status)
	GetDesiredStatus(actorKey Key) Status
	SetDesiredStatus(actorKey Key, status Status)
	GetInitialStatus(actorKey Key) Status
	GetTerminalStatus(actorKey Key) Status
}
