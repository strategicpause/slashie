package actor

type BasicActor struct {
	Type    Type
	Id      Id
	Mailbox chan func()
}

func (ba *BasicActor) Init() {
	for {
		select {
		case callback := <-ba.Mailbox:
			callback()
		}
	}
}

func (ba *BasicActor) GetType() Type {
	return ba.Type
}

func (ba *BasicActor) GetId() Id {
	return ba.Id
}

func (ba *BasicActor) Notify(callback func()) {
	ba.Mailbox <- callback
}
