package slashie

import (
	"fmt"
	"github.com/strategicpause/slashie/actor"
	"github.com/strategicpause/slashie/dependency"
	"github.com/strategicpause/slashie/logger"
	"github.com/strategicpause/slashie/subscription"
	"github.com/strategicpause/slashie/transition"
)

const (
	DefaultMailboxSize = 100
)

type slashie struct {
	actorRegistry       actor.Registry
	actorStatusManager  actor.StatusManager
	subscriptionManager subscription.Manager
	transitionManager   transition.Manager
	dependencyManager   dependency.Manager
	logger              logger.Logger
	mailbox             actor.Mailbox
}

type Opt func(s *slashie)

func WithMailboxSize(size int) Opt {
	return func(s *slashie) {
		s.mailbox = make(actor.Mailbox, size)
	}
}

func WithLogger(l logger.Logger) Opt {
	return func(s *slashie) {
		s.logger = l
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
	if s.subscriptionManager == nil {
		s.subscriptionManager = subscription.NewManager()
	}
	if s.transitionManager == nil {
		s.transitionManager = transition.NewManager()
	}
	if s.dependencyManager == nil {
		s.dependencyManager = dependency.NewManager()
	}
	if s.logger == nil {
		s.logger = logger.NewNullOutputLogger()
	}
	if s.mailbox == nil {
		s.mailbox = make(actor.Mailbox, DefaultMailboxSize)
	}

	go s.init()

	return s
}

func (s *slashie) init() {
	for msg := range s.mailbox {
		msg()
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

		actorKey := a.GetKey()
		if ok := s.actorRegistry.IsRegistered(a); !ok {
			errChan <- fmt.Errorf("unknown actor %s", actorKey)
			return
		}

		errChan <- s.updateStatus(actorKey, status)
	}
	return <-errChan
}

func (s *slashie) updateStatus(actorKey actor.Key, desiredStatus actor.Status) error {
	currentDesiredStatus := s.actorStatusManager.GetDesiredStatus(actorKey)
	if currentDesiredStatus == desiredStatus {
		s.logger.Debugf("%s desired status is already set to %s. Skipping update.", actorKey, desiredStatus)
		return nil
	}
	currentKnownStatus := s.actorStatusManager.GetKnownStatus(actorKey)
	if currentKnownStatus == desiredStatus {
		s.logger.Debugf("%s known status is already set to %s. Skipping update.", actorKey, desiredStatus)
		return nil
	}
	// Check to see if the current actor is already undergoing a transition. If so, then let's revisit this later.
	if currentDesiredStatus != currentKnownStatus {
		s.mailbox <- func() {
			s.logger.Debugf("%s is already transition from %s to %s. Deferring update.", actorKey, currentKnownStatus, currentDesiredStatus)
			if err := s.updateStatus(actorKey, desiredStatus); err != nil {
				s.logger.Debugf("%s", err)
			}
		}
		return nil
	}

	// Is this transition valid?
	ok := s.transitionManager.IsValidTransition(actorKey, currentDesiredStatus, desiredStatus)
	if !ok {
		return fmt.Errorf("transitioning from %s to %s is an illegal transition for actor %s", currentDesiredStatus, desiredStatus, actorKey)
	}

	s.logger.Debugf("Setting %s desired status to %s", actorKey, desiredStatus)
	s.actorStatusManager.SetDesiredStatus(actorKey, desiredStatus)

	s.mailbox <- func() {
		s.performTransition(actorKey)
	}

	return nil
}

func (s *slashie) performTransition(actorKey actor.Key) {
	desiredStatus := s.actorStatusManager.GetDesiredStatus(actorKey)
	// This will check to see if the current actor has any dependencies that it must wait for to transition.
	// If so, then this will block the current actor from transitioning to its desired status.
	hasDependencies := s.dependencyManager.HasTransitionDependencies(actorKey, desiredStatus)
	if hasDependencies {
		s.logger.Debugf("%s has a transition dependencies to %s", actorKey, desiredStatus)
		return
	}

	a, ok := s.actorRegistry.GetActor(actorKey)
	if !ok {
		s.logger.Debugf("could not find actor for %s", actorKey)
		return
	}

	knownStatus := s.actorStatusManager.GetKnownStatus(actorKey)
	// If we get this far, then we must first execute all actions before we can transition from the knownStatus
	// to the desiredStatus.
	s.logger.Debugf("Starting transition for %s: %s -> %s", actorKey, knownStatus, desiredStatus)
	s.transitionManager.StartTransition(actorKey, knownStatus, desiredStatus, func(action transition.Action) {
		a.Notify(func() {
			err := action()
			s.completeAction(actorKey, err)
		})
	})
}

func (s *slashie) completeAction(actorKey actor.Key, result error) {
	s.mailbox <- func() {
		s.transitionManager.CompleteTransitionAction(actorKey, result, func(results chan error) {
			newStatus := s.actorStatusManager.GetDesiredStatus(actorKey)
			for r := range results {
				if r != nil {
					s.logger.Debugf("%s", r)
					newStatus = s.actorStatusManager.GetTerminalStatus(actorKey)
					break
				}
			}
			s.updateKnownStatus(actorKey, newStatus)
		})
	}
}

func (s *slashie) updateKnownStatus(actorKey actor.Key, newStatus actor.Status) {
	// Execute any subscriptions that are waiting for the actor to transition.
	if a, ok := s.actorRegistry.GetActor(actorKey); ok {
		s.subscriptionManager.HandleSubscriptionsForStatus(actorKey, newStatus, func(s subscription.Subscription) {
			a.Notify(func() {
				s()
			})
		})
	}

	s.logger.Debugf("Setting known status for %s to %s.", actorKey, newStatus)
	s.actorStatusManager.SetKnownStatus(actorKey, newStatus)

	// Notify all dependencies that the current actor transitioned to the new status. This might result in other actors
	// transitioning to their destination status.
	s.dependencyManager.NotifyDependenciesOfStatus(actorKey, newStatus, func(depToNotify actor.Key) {
		s.mailbox <- func() {
			s.logger.Debugf("Notifying %s.", depToNotify)
			s.performTransition(depToNotify)
		}
	})

	terminalStatus := s.actorStatusManager.GetTerminalStatus(actorKey)
	if newStatus == terminalStatus {
		if a, ok := s.actorRegistry.GetActor(actorKey); ok {
			s.logger.Debugf("Stopping %s", actorKey)
			a.Stop()
		}
	}
}

func (s *slashie) AddTransitionDependency(srcActor actor.Actor, srcStatus actor.Status, depActor actor.Actor, depStatus actor.Status) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)

		srcKey := srcActor.GetKey()
		if ok := s.actorRegistry.IsRegistered(srcActor); !ok {
			errChan <- fmt.Errorf("unknown actor %s", srcKey)
			return
		}

		depKey := depActor.GetKey()
		if ok := s.actorRegistry.IsRegistered(depActor); !ok {
			errChan <- fmt.Errorf("unknown actor %s", depKey)
			return
		}

		errChan <- s.dependencyManager.AddTransitionDependency(srcKey, srcStatus, depKey, depStatus)
	}
	return <-errChan
}

func (s *slashie) AddTransitionActions(actor actor.Actor, actions []*transition.TransitionAction) error {
	for _, action := range actions {
		if err := s.AddTransitionAction(actor, action.SrcStatus, action.DestStatus, action.Action); err != nil {
			return err
		}
	}
	return nil
}

func (s *slashie) AddTransitionAction(a actor.Actor, srcStatus actor.Status, destStatus actor.Status, action transition.Action) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)
		actorKey := a.GetKey()

		if isValid := s.actorStatusManager.IsValidTransitionStatus(actorKey, srcStatus, destStatus); !isValid {
			errChan <- fmt.Errorf("cannot transition from %s to %s", srcStatus, destStatus)
			return
		}
		s.logger.Debugf("Adding transaction action for %s for %s -> %s.", actorKey, srcStatus, destStatus)
		s.transitionManager.AddTransitionAction(actorKey, srcStatus, destStatus, action)
	}
	return <-errChan
}

func (s *slashie) GetStatus(a actor.Actor) actor.Status {
	responseChan := make(chan actor.Status)
	s.mailbox <- func() {
		defer close(responseChan)

		actorKey := a.GetKey()
		responseChan <- s.actorStatusManager.GetKnownStatus(actorKey)
	}
	return <-responseChan
}

func (s *slashie) Subscribe(a actor.Actor, status actor.Status, callback subscription.Subscription) error {
	errChan := make(chan error)
	s.mailbox <- func() {
		defer close(errChan)

		actorKey := a.GetKey()
		isValid := s.actorStatusManager.IsValidSubscriptionStatus(actorKey, status)
		if !isValid {
			errChan <- fmt.Errorf("cannot subscribe to current status %s", status)
			return
		}

		s.subscriptionManager.Subscribe(actorKey, status, callback)
	}
	return <-errChan
}
