package subscription

import "github.com/strategicpause/slashie/actor"

type manager struct {
	// subscriptionsForActor
	subscriptionsForActor map[actor.Key]SubscriptionsByStatus
}

func NewManager() Manager {
	return &manager{
		subscriptionsForActor: map[actor.Key]SubscriptionsByStatus{},
	}
}

func (m *manager) Subscribe(actorKey actor.Key, status actor.Status, callback Subscription) {
	if _, ok := m.subscriptionsForActor[actorKey]; !ok {
		m.subscriptionsForActor[actorKey] = SubscriptionsByStatus{}
	}
	subscriptionsByStatus := m.subscriptionsForActor[actorKey]

	subscriptionsByStatus[status] = append(subscriptionsByStatus[status], callback)
}

func (m *manager) HandleSubscriptionsForStatus(actorKey actor.Key, status actor.Status, callback func(s Subscription)) {
	if _, ok := m.subscriptionsForActor[actorKey]; !ok {
		return
	}
	subscriptions := m.subscriptionsForActor[actorKey][status]
	for _, subscription := range subscriptions {
		callback(subscription)
	}
	delete(m.subscriptionsForActor[actorKey], status)
}
