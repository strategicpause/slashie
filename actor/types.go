package actor

// Key is a composite of the Type and Id which will uniquely identify an actor.
type Key string

// Type specifies the type of actor.
type Type string

// Id uniquely identifies an actor of a given type.
type Id string

type Status string

type Message func()

type Mailbox chan Message
