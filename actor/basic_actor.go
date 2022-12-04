package actor

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	DefaultMailBoxSize = 100
)

type BasicActor struct {
	actorType Type
	actorId   Id
	mailbox   mailbox
	stopChan  chan bool
	wg        sync.WaitGroup
	handlers  map[reflect.Type]Handler
}

func NewBasicActor(actorType Type, actorId Id) *BasicActor {
	ba := &BasicActor{
		actorType: actorType,
		actorId:   actorId,
		mailbox:   make(mailbox, DefaultMailBoxSize),
		stopChan:  make(chan bool, 1),
		wg:        sync.WaitGroup{},
		handlers:  map[reflect.Type]Handler{},
	}
	// Bootstrap message handlers
	ba.registerMessageHandler(messageType, ba.handleMessage)

	ba.wg.Add(1)
	go ba.Init()

	return ba
}

func (ba *BasicActor) handleMessage(message any) {
	message.(Message)()
}

func (ba *BasicActor) Init() {
	for {
		select {
		case message := <-ba.mailbox:
			messageType := reflect.TypeOf(message)
			if handler, ok := ba.handlers[messageType]; ok {
				handler(message)
			}
		case <-ba.stopChan:
			return
		}
	}
}

func (ba *BasicActor) GetType() Type {
	return ba.actorType
}

func (ba *BasicActor) GetId() Id {
	return ba.actorId
}

func (ba *BasicActor) GetKey() Key {
	actorKey := fmt.Sprintf("%s:%s", ba.GetType(), ba.GetId())
	return Key(actorKey)
}

func (ba *BasicActor) Notify(message Message) {
	ba.mailbox <- message
}

func (ba *BasicActor) RegisterMessageHandler(messageType any, handler Handler) {
	ba.mailbox <- Message(func() {
		ba.registerMessageHandler(messageType, handler)
	})
}

func (ba *BasicActor) registerMessageHandler(messageType any, handler Handler) {
	reflectType := reflect.TypeOf(messageType)
	ba.handlers[reflectType] = handler
}

func (ba *BasicActor) SendMessage(message any) error {
	errChan := make(chan error)
	ba.mailbox <- Message(func() {
		defer close(errChan)

		messageType := reflect.TypeOf(message)
		if _, ok := ba.handlers[messageType]; !ok {
			errChan <- fmt.Errorf("unknown message type: %s", messageType)
		}

	})

	if err := <-errChan; err != nil {
		return err
	}
	ba.mailbox <- message

	return nil
}

func (ba *BasicActor) Stop() {
	ba.mailbox <- Message(func() {
		ba.stopChan <- true
		ba.wg.Done()
	})
}

func (ba *BasicActor) Wait() {
	ba.wg.Wait()
}
