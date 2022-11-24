package actor

// Actor defines the operations expected that each actor implements.
type Actor interface {
	// GetType returns the type of Actor.
	GetType() Type
	// GetId returns the Id for an Actor. It is expected that Id is unique for a given Type of Actor.
	GetId() Id
	// GetKey will return the key for the given Actor
	GetKey() Key
	// Notify will send a message
	Notify(message Message)
	// Init
	Init()
	// Stop
	Stop()
}

type Registry interface {
	// RegisterActor will register an Actor and return the corresponding ActorKey.
	RegisterActor(actor Actor) Key
	// GetActor will return an Actor, if it exists, given an ActorKey.
	GetActor(actorKey Key) (Actor, bool)
	//
	IsRegistered(actor Actor) bool
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
