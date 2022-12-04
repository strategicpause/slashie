package actor

import (
	"github.com/stretchr/testify/assert"
	"sync"
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

type testType struct {
	message string
}

func TestBasicActor_SendMessage(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	wg := sync.WaitGroup{}
	wg.Add(1)

	messageHandled := false
	testMessage := "TestMessage"
	actor.RegisterMessageHandler(testType{}, func(message any) {
		messageHandled = true
		assert.Equal(t, testMessage, message.(testType).message)

		wg.Done()
	})

	err := actor.SendMessage(testType{message: testMessage})
	assert.Nil(t, err)

	wg.Wait()
	assert.True(t, messageHandled)
}

func TestBasicActor_SendUnknownMessage(t *testing.T) {
	actor := NewBasicActor(ActorType, ActorId)

	err := actor.SendMessage(testType{message: "TestMessage"})

	assert.NotNil(t, err)
}
