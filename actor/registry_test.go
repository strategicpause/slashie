package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ActorType = Type("Actor")
	ActorId   = Id("ActorA")
	ActorKey  = Key("Actor-ActorA")
)

func TestRegisterActor(t *testing.T) {
	registry := NewRegistry()
	actor := NewBasicActor(ActorType, ActorId)

	key := registry.RegisterActor(actor)

	assert.Equal(t, ActorKey, key)
}

func NewBasicActor(actorType Type, actorId Id) Actor {
	return &BasicActor{
		Type: actorType,
		Id:   actorId,
	}
}

func TestGetActorKey_Registered(t *testing.T) {
	registry := NewRegistry()
	actor := NewBasicActor(ActorType, ActorId)
	actorKey := registry.RegisterActor(actor)

	newActorKey, ok := registry.GetActorKey(actor)

	assert.Equal(t, actorKey, newActorKey)
	assert.True(t, ok)
}

func TestGetActorKey_NotRegistered(t *testing.T) {
	registry := NewRegistry()
	actor := NewBasicActor(ActorType, ActorId)

	actorKey, ok := registry.GetActorKey(actor)
	assert.Equal(t, Key(""), actorKey)
	assert.False(t, ok)
}
