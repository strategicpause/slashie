package actor

import (
	"fmt"
	"sync"
)

const (
	DefaultMailBoxSize = 100
)

type BasicActor struct {
	actorType Type
	actorId   Id
	mailbox   Mailbox
	stopChan  chan bool
	wg        sync.WaitGroup
}

func NewBasicActor(actorType Type, actorId Id) *BasicActor {
	ba := &BasicActor{
		actorType: actorType,
		actorId:   actorId,
		mailbox:   make(Mailbox, DefaultMailBoxSize),
		stopChan:  make(chan bool, 1),
		wg:        sync.WaitGroup{},
	}

	ba.wg.Add(1)
	go ba.Init()

	return ba
}

func (ba *BasicActor) Init() {
	for {
		select {
		case message := <-ba.mailbox:
			message()
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

func (ba *BasicActor) Stop() {
	ba.mailbox <- func() {
		ba.stopChan <- true
		ba.wg.Done()
	}
}

func (ba *BasicActor) Wait() {
	ba.wg.Wait()
}
