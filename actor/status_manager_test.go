package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	InitStatus     = Status("Init")
	MidStatus      = Status("Mid")
	TerminalStatus = Status("Terminal")
)

func TestInitializeActor(t *testing.T) {
	mgr := NewStatusManager()
	mgr.InitializeActor(ActorKey, InitStatus, TerminalStatus)

	assert.Equal(t, InitStatus, mgr.GetInitialStatus(ActorKey))
	assert.Equal(t, TerminalStatus, mgr.GetTerminalStatus(ActorKey))
	assert.Equal(t, InitStatus, mgr.GetDesiredStatus(ActorKey))
	assert.Equal(t, InitStatus, mgr.GetKnownStatus(ActorKey))
}

func TestIsValidSubscriptionStatus(t *testing.T) {
	mgr := NewStatusManager()
	mgr.InitializeActor(ActorKey, InitStatus, TerminalStatus)
	mgr.SetKnownStatus(ActorKey, MidStatus)

	// A previous status is not valid
	assert.False(t, mgr.IsValidSubscriptionStatus(ActorKey, InitStatus))
	// The current status is not valid
	assert.False(t, mgr.IsValidSubscriptionStatus(ActorKey, MidStatus))
	// A non-visited status is valid
	assert.True(t, mgr.IsValidSubscriptionStatus(ActorKey, TerminalStatus))
}

func TestIsValidTransitionStatus(t *testing.T) {
	mgr := NewStatusManager()
	mgr.InitializeActor(ActorKey, InitStatus, TerminalStatus)
	mgr.SetKnownStatus(ActorKey, MidStatus)

	assert.False(t, mgr.IsValidTransitionStatus(ActorKey, TerminalStatus, MidStatus))
	assert.False(t, mgr.IsValidTransitionStatus(ActorKey, MidStatus, MidStatus))
	assert.False(t, mgr.IsValidTransitionStatus(ActorKey, MidStatus, InitStatus))
	assert.True(t, mgr.IsValidTransitionStatus(ActorKey, InitStatus, MidStatus))
	assert.True(t, mgr.IsValidTransitionStatus(ActorKey, MidStatus, TerminalStatus))
	assert.True(t, mgr.IsValidTransitionStatus(ActorKey, InitStatus, TerminalStatus))
}