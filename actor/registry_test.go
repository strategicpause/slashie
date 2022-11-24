package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ActorType = Type("Actor")
	ActorId   = Id("ActorA")
	ActorKey  = Key("Actor:ActorA")
)

func TestRegisterActor(t *testing.T) {
	registry := NewRegistry()
	actor := NewBasicActor(ActorType, ActorId)

	key := registry.RegisterActor(actor)

	assert.Equal(t, ActorKey, key)
}

func TestGetActor(t *testing.T) {
	registry := NewRegistry()
	actor := NewBasicActor(ActorType, ActorId)
	actorKey := registry.RegisterActor(actor)

	a, ok := registry.GetActor(actorKey)
	assert.True(t, ok)
	assert.Equal(t, actor, a)
}

func TestGetActor_NotRegistered(t *testing.T) {
	registry := NewRegistry()

	a, ok := registry.GetActor("actor-Key")
	assert.False(t, ok)
	assert.Nil(t, a)
}
