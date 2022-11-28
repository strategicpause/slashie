# Slashie
Slashie is an implementation of the actor model in Go in which a goroutine represents an actor. 
Slashie allows users to configure transition rules for actors, including actions that must run before an actor can 
transition to a desired state, and whether or not any other actors must first be in a given state before transitioning. 
State management is centralized such that the actor graph and transition rules can be validated. Slashie is also the 
name of the highly prestigious award given at the VH1 Fashion Awards to actor "slash" models (and not the other way 
around).

# Tenets
* **Testable first**. All code and functionality should be covered by unit tests to validate behavior. A feature is not
considered done unless it has tests to validate behavior.
* **Minimal Footprint & minimal dependencies**. Slashie will carefully evaluate the trade-offs of bringing in a new 
dependency versus implementing it from scratch. Each dependency added becomes a dependency of all software packages
which utilize slashie.
* **Centralize Actor State**. Rather than having each actor maintain its own state, this information is instead 
centralized such that the entire actor graph can be validated and optimal decisions can be made.

# Actor Model
An actor can make local decisions (ie: update local state), or send messages to other actors. One actor cannot mutate
the state of another actor. Instead actors communicate by sending messages which may result in an actor taking some 
action which may result in mutating its own local state.

When an actor receives a new message, it is queued in the mailbox which is represented as a Go buffered channel. As a
result each actor will handle messages in the order that they are received. As a result actors don't have to use locks
when mutating state.

https://doc.akka.io/docs/akka/current/typed/guide/actors-intro.html
https://en.wikipedia.org/wiki/Actor_model
https://www.brianstorti.com/the-actor-model/

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
- Logging - Know when an actor wants to transition, has transitioned, is notifying dependencies, etc.
- Message Passing - An actor should be able to pass a message to an arbitrary actor without having a reference to that actor.
- Actor crash management - Can we recover? State persistence.
- Expand to processes. Can we decouple actors from goroutines and extend the definition to processes? What about a process on a separate machine? 
