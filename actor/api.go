package actor

// Actor defines the operations that each actor implements.
type Actor interface {
	// GetType returns the type of Actor.
	GetType() Type
	// GetId returns the Id for an Actor. It is expected that Id is unique for a given Type of Actor.
	GetId() Id
	// GetKey will return the key for the given Actor
	GetKey() Key
	// Notify will send a message
	Notify(message Message)
	// Init will initialize the event loop for handling messages that get sent to the actor's mailbox.
	Init()
	// Stop will stop all event processing and kill the underlying goroutine.
	Stop()
}

// Registry is a central repository to register & fetch actors.
type Registry interface {
	// RegisterActor will register an Actor and return the corresponding actor Key.
	RegisterActor(actor Actor) Key
	// GetActor will return an Actor, if it exists, given an ActorKey. The second parameter will return false if an
	// actor does not exist for the given Key.
	GetActor(actorKey Key) (Actor, bool)
	// IsRegistered returns true if the given actor was registered with via RegisterActor.
	IsRegistered(actor Actor) bool
}

// StatusManager keeps track of actor statues including the inital, terminal, desired, and known status.
type StatusManager interface {
	// InitializeActor will set the initial and terminal status for the given actor. The desired and known status
	// will both be set to the initial status.
	InitializeActor(actorKey Key, initStatus Status, terminalStatus Status)
	// IsValidSubscriptionStatus returns true if the actor can transition to the given status. You cannot
	// subscribe to a status which has already past.
	IsValidSubscriptionStatus(actorKey Key, status Status) bool
	// IsValidTransitionStatus returns false for illegal transitions including using the terminal status as the source
	// or using the init status as the destination.
	IsValidTransitionStatus(actorKey Key, srcStatus Status, destStatus Status) bool
	// GetKnownStatus returns the known status for the given actor Key.
	GetKnownStatus(actorKey Key) Status
	// SetKnownStatus sets the known status for the given actor Key.
	SetKnownStatus(actorKey Key, knownStatus Status)
	// GetDesiredStatus returns the desired status for the given actor Key.
	GetDesiredStatus(actorKey Key) Status
	// SetDesiredStatus sets the desired status for the given actor Key.
	SetDesiredStatus(actorKey Key, status Status)
	// GetInitialStatus returns the initial status for the given actor Key.
	GetInitialStatus(actorKey Key) Status
	// GetTerminalStatus returns the terminal status for the given actor Key.
	GetTerminalStatus(actorKey Key) Status
}
