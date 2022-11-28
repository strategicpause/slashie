package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetType(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	assert.Equal(t, ActorType, actor.GetType())
}

func TestGetId(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	assert.Equal(t, ActorId, actor.GetId())
}

func TestGetKey(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	assert.Equal(t, ActorKey, actor.GetKey())
}

func TestNotify(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	messageProcessed := false

	actor.Notify(func() {
		messageProcessed = true
	})
	actor.Stop()
	actor.Wait()

	assert.True(t, messageProcessed)
}
