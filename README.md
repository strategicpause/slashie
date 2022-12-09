[![Go Report Card](https://goreportcard.com/badge/github.com/strategicpause/slashie)](https://goreportcard.com/report/github.com/strategicpause/slashie)

# Slashie
Slashie is an implementation of the actor model in Go in which a goroutine represents an actor. Slashie is also a 
state machine library which also manages actor status and managing from one status to the next. Slashie allows users to 
configure transition rules for actors, including actions that must run before an actor can transition to a desired 
status, and whether or not any other actors must first be in a given status before transitioning. Status management is 
centralized such that the actor graph and transition rules can be validated. 

Slashie is also the name of the highly prestigious award given at the VH1 Fashion Awards to actor "slash" models (and 
not the other way around).

## Features
* **Transitions**: An actor may behave differently based on its current status. Slashie allows users to define
valid status transitions, and actions to take in order to reach that status.
* **Dependencies**: In order for an actor to transition to a given status or status, it may depend on another actor
transitioning to a given status. For example, maybe in order for an actor to communicate with a remote entity, networking
must be setup first. Slashie allows users to establish dependencies before an actor is eligible to transition.
* **Subscriptions**: An actor can subscribe a given actor transitioning to a given status. This results the subscribing
actor receiving a message with some action to take.
* **Centralized state**: Centralizing status & transition information allows slashie to make smart & efficient decisions.
For example slashie can detect when adding a new dependency will introduce circular dependencies. Slashie will also be 
able to notify an actor once all dependencies have transition to their target status.
* **Message Passing**: Actors communicate by passing messages to one another. Slashie provides a mechanism to allow an
arbitrary actor to send a message to another actor without requiring a direct reference to them.
* **Tested**: Slashie has over 90% test coverage, including both unit & integration test coverage. This can provide
immediate feedback to developers as to whether or not their feature is introducing any bugs or affects any behavior
that users expect.

## Tenets
* **Testable first**. All code and functionality should will be covered by unit tests to validate behavior. A feature is 
not considered done unless it has tests which demonstrate that the feature works.
* **Minimal Footprint & minimal dependencies**. Slashie will favor not introducing a new dependency since any new 
dependency now becomes a dependency of all software packages which utilize slashie. Functionality which is coupled to
a specific dependency will live in a separate package from the core.
* **Centralize Actor State**. Rather than having each actor maintain its own status, this information is instead 
centralized such that the entire actor graph can be validated and optimal decisions can be made.

# Actor Model
The actor model uses the actor as the basic primitive of concurrent computation. At the heart of an actor is an event
loop which reads messages from a buffered channel which is referred to as a mailbox. When a message is read from a 
mailbox, the actor is actually executing some function which could represent business logic that makes decisions (ie: 
updates local state), performs some action in a larger workflow, sends notifications to other actors, or perhaps spawns 
other actors. 

One actor cannot mutate the state of another actor. Instead actors communicate by sending messages to each other. When 
an actor receives a new message, it is queued. As a result 
each actor will handle exactly one message at a time in the order that they are received. By only handling one message 
at a time, actors don't have to use synchronization mechanisms like locks or semaphores when mutating state.

Slashie represents an actor as a Goroutine which is either waiting for a message from the mailbox or executing a message. 
When a message is received it will execute that function. The function could be a wrapper around a transition action,
a subscription, or perhaps an arbitrary message which originates from another actor. The following sections talk more
about features of slashie.

# State Machine
TODO

## Status & Actions

## Dependencies

## Subscriptions

# Usage
~~~~
$ cat main.go

package main

const (
	InitStatus  = "Init"
	PrintStatus = "Print"
)

func main() {
    # Initialize Slashie.
    slashie := slashie.NewSlashie()
    # Initialize a basic actor.
    actor := actor.NewBasicActor("ExampleActor", "HelloWorld")
    # Register the actor with Slashie. An initial and terminal status must be specified.
    slashie.AddActor(actor, InitStatus, PrintStatus)
    // Register an action which will print "Hello " before the actor trainsitions from Init to Print.
    // By having the registered action return a nil error, the actor will transition to the PrintStatus. 
    slashie.AddTransitionAction(actor, InitStatus, PrintStatus, func() error {
        fmt.Printf("Hello ")
        return nil
    })
    // Print "World!" after the actor transitions to Done.
    slashie.Subscribe(actor, PrintStatus, func() {
        fmt.Println("World!")
    })
    // Indicate to slashie to transition the actor to transition to the "Print" status.
    slashie.UpdateStatus(actor, PrintStatus)
    // Block until the actor has transitioned to the terminal status.
    actor.Wait()
}

$ go run main.go
Hello World!
~~~~
Additional examples:
* [Weedupe](https://github.com/strategicpause/slashie-weedupe) - A mini map-reduce library, inspired by Hadoop, written
using Slashie. A director actor will spin up multiple actors to count the number of times a word appears in a given set 
of files. The results will be combined and printed to the screen.
# Testing
~~~~
# Run all tests
make test

# Validate code coverage
make coverage

# View code coverage report
make coverage-report
~~~~

# TODO
- Actor crash management - Can we recover? State persistence. Provide interface in this package, but keep implementation in separate packages. This would allow us to easily switch implementation for different dbs such as etcd or BoltDb.
- Expand to processes. Can we decouple actors from goroutines and extend the definition to processes? What about a process on a separate machine? 
- Export to SVG. The ability to export the transition graph + dependencies to a visual representation.
- Metrics observers. ie: Observer which runs every N seconds that provides mailbox size. Provide interface in main package, but implementation in separate package.
  - Time message spent in the mailbox queue. This will help tune mailbox sizes. A smaller mailbox may mean that
messages don't received by an actor. A larger mailbox means that messages may wait longer in the mailbox. Measuring this will help determine what is acceptable.
- Supervisor Tree - Have a parent get notified if a child actor fails. That way the parent can take some action like spin 
up a new actor or change its own state.