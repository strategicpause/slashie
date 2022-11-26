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

// Mailbox is used to store incoming messages which are waiting to be processed.
type Mailbox chan Message
