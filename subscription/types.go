package subscription

import "github.com/strategicpause/slashie/actor"

// Subscription is a function to execute after an actor transitions to a given state.
type Subscription func()

// SubscriptionsByStatus
type SubscriptionsByStatus map[actor.Status][]Subscription
