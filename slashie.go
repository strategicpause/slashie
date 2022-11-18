package slashie

import (
	"errors"
	"fmt"
	"github.com/strategicpause/slashie/actor"
	"github.com/strategicpause/slashie/transition"
	"sync"
)

const (
	UnknownStatus      = actor.Status("UNKNOWN")
	DefaultMailboxSize = 100
)

type slashie struct {
	actorRegistry      actor.Registry
	actorStatusManager actor.StatusManager
	transitionManager  transition.Manager
	mailbox            chan func()
}

type Opt func(s *slashie)

func WithActorRegistry(registry actor.Registry) Opt {
	return func(s *slashie) {
		s.actorRegistry = registry
	}
}

func WithActorStatusManager(manager actor.StatusManager) Opt {
	return func(s *slashie) {
		s.actorStatusManager = manager
	}
}

func WithTransitionManager(manager transition.Manager) Opt {
	return func(s *slashie) {
		s.transitionManager = manager
	}
}

func WithMailboxSize(size int) Opt {
	return func(s *slashie) {
		s.mailbox = make(chan func(), size)
	}
}

func NewSlashie(opts ...Opt) Slashie {
	s := &slashie{}

	for _, opt := range opts {
		opt(s)
	}

	// Set defaults
	if s.actorRegistry == nil {
		s.actorRegistry = actor.NewRegistry()
	}
	if s.actorStatusManager == nil {
		s.actorStatusManager = actor.NewStatusManager()
	}
	if s.transitionManager == nil {
		s.transitionManager = transition.NewManager()
	}
	if s.mailbox == nil {
		s.mailbox = make(chan func(), DefaultMailboxSize)
	}

	go s.init()

	return s
}

func (s *slashie) init() {
	for {
		select {
		case msg := <-s.mailbox:
			msg()
		}
	}
}

func (s *slashie) AddActor(actor actor.Actor, initStatus actor.Status, terminalStatus actor.Status) {
	s.mailbox <- func() {
		actorKey := s.actorRegistry.RegisterActor(actor)
		s.actorStatusManager.InitializeActor(actorKey, initStatus, terminalStatus)
	}
}

func (s *slashie) UpdateStatus(a actor.Actor, status actor.Status) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)
		actorKey, ok := s.actorRegistry.GetActorKey(a)
		if !ok {
			errChan <- errors.New("unknown actor")
		}
		s.updateStatus(actorKey, status)
	}
	return <-errChan
}

func (s *slashie) updateStatus(actorKey actor.Key, desiredStatus actor.Status) {
	// Check to see if this is a legal transition
	s.actorStatusManager.SetDesiredStatus(actorKey, desiredStatus)

	// This will check to see if the current actor has any dependencies that it must wait for to transition.
	// If so, then this will block the current actor from transitioning to its desired status.
	hasDependencies := s.transitionManager.CanTransitionToStatus(actorKey, desiredStatus)
	if hasDependencies {
		return
	}

	// If we get this far, then we must first execute all callbacks before we can transition from the currentStatus
	// to the desiredStatus.
	currentStatus := s.actorStatusManager.GetKnownStatus(actorKey)
	callbacks := s.transitionManager.GetTransitionCallbacks(actorKey, currentStatus, desiredStatus)
	a := s.actorRegistry.GetActor(actorKey)

	// Prepare a channel which will capture any errors returned by the callback functions.
	numCallbacks := len(callbacks)
	errs := make(chan error, numCallbacks)
	// Prepare a WaitGroup which will block until all the callback functions are done executing.
	wg := sync.WaitGroup{}
	wg.Add(numCallbacks)

	// Execute callback functions by sending it to the current actor's mailbox (via Notify).
	for _, callback := range callbacks {
		cb := func() {
			errs <- callback()
			wg.Done()
		}
		a.Notify(cb)
	}

	// This will wait until all callbacks are done executing before updating the known status
	go func() {
		wg.Wait()
		close(errs)
		s.mailbox <- func() {
			// If there were no errors from the transition callbacks, then we're able to transition to the
			// desiredStatus. Otherwise, we will transition to the terminalStatus.
			newStatus := desiredStatus
			for err := range errs {
				if err != nil {
					newStatus = s.actorStatusManager.GetTerminalStatus(actorKey)
					break
				}
			}
			s.updateKnownStatus(actorKey, newStatus)
		}
	}()
}

func (s *slashie) updateKnownStatus(actorKey actor.Key, newStatus actor.Status) {
	s.actorStatusManager.SetKnownStatus(actorKey, newStatus)
	// Now that the given actor has transitioned to the new status, we can execute any subscriptions that are waiting
	// for the actor to transition.
	subscriptions := s.transitionManager.GetSubscriptionsForStatus(actorKey, newStatus)
	numSubscriptions := len(subscriptions)
	if numSubscriptions > 0 {
		for _, subscription := range subscriptions {
			subscription()
		}
	}
	// Notify all dependencies that the current actor transitioned to the new status. This might result in other actors
	// transitioning to their destination status.
	s.transitionManager.NotifyDependenciesOfStatus(actorKey, newStatus, func(depKey actor.Key, depStatus actor.Status) {
		s.mailbox <- func() {
			s.updateStatus(depKey, depStatus)
		}
	})
}

func (s *slashie) AddTransitionDependency(srcActor actor.Actor, srcStatus actor.Status, depActor actor.Actor, depStatus actor.Status) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)
		srcKey, ok := s.actorRegistry.GetActorKey(srcActor)
		if !ok {
			errChan <- errors.New("unknown source actor")
			return
		}

		depKey, ok := s.actorRegistry.GetActorKey(depActor)
		if !ok {
			errChan <- errors.New("unknown dependent actor")
			return
		}

		errChan <- s.transitionManager.AddTransitionDependency(srcKey, srcStatus, depKey, depStatus)
	}
	return <-errChan
}

func (s *slashie) AddTransitionCallbacks(actor actor.Actor, tuples []*transition.CalbackTuple) error {
	for _, tuple := range tuples {
		err := s.AddTransitionCallback(actor, tuple.SrcStatus, tuple.DestStatus, tuple.Callback)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *slashie) AddTransitionCallback(actor actor.Actor, srcStatus actor.Status, destStatus actor.Status, callback transition.Callback) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)
		actorKey, ok := s.actorRegistry.GetActorKey(actor)
		if !ok {
			errChan <- errors.New("unknown actor")
			return
		}

		isValid := s.actorStatusManager.IsValidTransitionStatus(actorKey, srcStatus, destStatus)
		if !isValid {
			errChan <- fmt.Errorf("cannot transition from %s to %s", srcStatus, destStatus)
			return
		}
		s.transitionManager.AddTransitionCallback(actorKey, srcStatus, destStatus, callback)
	}
	return <-errChan
}

func (s *slashie) GetStatus(a actor.Actor) actor.Status {
	responseChan := make(chan actor.Status)
	s.mailbox <- func() {
		defer close(responseChan)
		actorKey, ok := s.actorRegistry.GetActorKey(a)
		if !ok {
			responseChan <- UnknownStatus
			return
		}
		responseChan <- s.actorStatusManager.GetKnownStatus(actorKey)
	}
	return <-responseChan
}

func (s *slashie) Subscribe(actor actor.Actor, status actor.Status, callback transition.Subscription) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)
		actorKey, ok := s.actorRegistry.GetActorKey(actor)
		if !ok {
			errChan <- errors.New("unknown actor")
			return
		}

		isValid := s.actorStatusManager.IsValidSubscriptionStatus(actorKey, status)
		if !isValid {
			errChan <- fmt.Errorf("cannot subscribe to current status %s", status)
			return
		}

		s.transitionManager.Subscribe(actorKey, status, callback)
	}
	return <-errChan
}
