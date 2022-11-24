package actor

type statusManager struct {
	// initialStatusByActor determines the initial default status for a given actor when it has not yet been
	// initialized. An actor cannot transition back to the initial status.
	initialStatusByActor map[Key]Status
	// terminalStatusByActor answers what status determines if an actor has reached a terminal state. Once in the
	// terminal state, an Actor cannot transition to another state.
	terminalStatusByActor map[Key]Status
	// desiredStatusByActor determines what status an Actor is attempting to transition to.
	desiredStatusByActor map[Key]Status
	// knownStatusByActor determines the current status of an Actor.
	knownStatusByActor map[Key]Status
}

func NewStatusManager() StatusManager {
	return &statusManager{
		initialStatusByActor:  map[Key]Status{},
		terminalStatusByActor: map[Key]Status{},
		desiredStatusByActor:  map[Key]Status{},
		knownStatusByActor:    map[Key]Status{},
	}
}

func (a *statusManager) InitializeActor(actorKey Key, initStatus Status, terminalStatus Status) {
	a.initialStatusByActor[actorKey] = initStatus
	a.terminalStatusByActor[actorKey] = terminalStatus

	a.desiredStatusByActor[actorKey] = initStatus
	a.knownStatusByActor[actorKey] = initStatus
}

func (a *statusManager) IsValidSubscriptionStatus(actorKey Key, status Status) bool {
	currentKnownStatus := a.knownStatusByActor[actorKey]
	if currentKnownStatus == status {
		return false
	}
	initialStatus := a.initialStatusByActor[actorKey]

	return initialStatus != status
}

func (a *statusManager) IsValidTransitionStatus(actorKey Key, srcStatus Status, destStatus Status) bool {
	if srcStatus == destStatus {
		return false
	}
	if srcStatus == a.terminalStatusByActor[actorKey] {
		return false
	}
	if destStatus == a.initialStatusByActor[actorKey] {
		return false
	}
	return true
}

func (a *statusManager) GetKnownStatus(actorKey Key) Status {
	return a.knownStatusByActor[actorKey]
}

func (a *statusManager) SetKnownStatus(actorKey Key, status Status) {
	a.knownStatusByActor[actorKey] = status
}

func (a *statusManager) GetDesiredStatus(actorKey Key) Status {
	return a.desiredStatusByActor[actorKey]
}

func (a *statusManager) SetDesiredStatus(actorKey Key, status Status) {
	a.desiredStatusByActor[actorKey] = status
}

func (a *statusManager) GetInitialStatus(actorKey Key) Status {
	return a.initialStatusByActor[actorKey]
}

func (a *statusManager) GetTerminalStatus(actorKey Key) Status {
	return a.terminalStatusByActor[actorKey]
}
