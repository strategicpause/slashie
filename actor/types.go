package actor

// Key is a composite of the Type and Id which will uniquely identify an actor.
type Key string

// Type specifies the type of actor.
type Type string

// Id uniquely identifies an actor of a given type.
type Id string

// Status is used to indicate state of the actor.
type Status string

// Message represents some unit of computation that the actor processes.
type Message func()

// messageType is used to register the Handler for the Message type.
var messageType = Message(func() {})

// mailbox is the internal mailbox type used to store incoming messages which are waiting to be processed.
type mailbox chan any

// Handler is a function which can process an incoming message to an actor.
type Handler func(message any)
