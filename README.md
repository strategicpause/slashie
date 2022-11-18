# Slashie
Slashie is an implementation of the actor model in Go which allows users to set transition rules for actors, and 
centralize state management. Slashie is also the name of the highly prestigious award given at the VH1 Fashion Awards 
to actor "slash" models (and not the other way around).

# Principles
* **Testable first**. All code and functionality should be covered by unit tests to validate behavior. If functionality cannot be tested, then it may indicate the wrong abstractions are being used or tight coupling between components.
* **Minimal Footprint & dependencies**. Slashie will carefully evaluate the trade-offs of bringing in a new dependency versus implementing it from scratch.    
* 

# Actor Model
TODO - What is the actor model?
* Instead of calling methods, actors send messages to one another.
* An actor can create other actors.
* Receive messages which result in performing some action such as mutate local state or perhaps send messages to other actors.
* An actor processes messages one message at a time. Messages are stored in a queue called a mailbox.
* An Actor can be a goroutine, a process on the same machine, or a process on a remote machine.
* Instead of a mailbox being in memory, could a mailbox be a queue on disk?
https://doc.akka.io/docs/akka/current/typed/guide/actors-intro.html
https://en.wikipedia.org/wiki/Actor_model
https://www.brianstorti.com/the-actor-model/

# Usage
~~~~
TODO - Basic usage
~~~~

# Testing
~~~~
make test
~~~~

# TODO
- Logging - Know when an actor wants to transition, has transitioned, is notifying dependencies, etc.
- Actor crash management - Can we recover? State persistence.
- Expand to processes. Can we decouple actors from goroutines and extend the defintiion to processes? What about a process on a separate machine? 
