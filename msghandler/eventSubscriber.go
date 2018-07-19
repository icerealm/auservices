package msghandler

import (
	"auservices/domain"
	"auservices/utilities"
	"log"
	"os"

	"github.com/nats-io/go-nats-streaming"
)

const (
	durableCategoryID = "durable-category-id"
	durableItemID     = "durable-item-id"
)

var eventSubscriber *EventSubscriber

//EventSubscriber collect subscriber list.
type EventSubscriber struct {
	subscriberMap map[string]*MessageSubscriber
	Subscribers   []*MessageSubscriber
}

//AppendSubscriber append a subscriber.
func (s *EventSubscriber) AppendSubscriber(channelID string, ms *MessageSubscriber) {
	if s.subscriberMap == nil {
		s.subscriberMap = make(map[string]*MessageSubscriber)
	}
	s.subscriberMap[channelID] = ms
	s.Subscribers = append(s.Subscribers, ms)
}

//GetEventSubscribers get eventsubscribers instance.
func GetEventSubscribers() *EventSubscriber {
	eventSubscriber = &EventSubscriber{}
	//initial category subscriber
	categorySub := initCategorySubscriber()
	eventSubscriber.AppendSubscriber(kcategoryChannelID, categorySub)

	return eventSubscriber
}

func initCategorySubscriber() *MessageSubscriber {
	//initial category subscriber
	sub, err := CreateMessageSubscriber(
		MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial category subscriber")
		os.Exit(1)
	}
	sub.SubscribeEvent(kcategoryChannelID, durableCategoryID, categoryEventMsgHandler)
	return sub
}

//categoryEventMsgHandler to handle business logic for category event
func categoryEventMsgHandler(msg *stan.Msg) {
	go func(m *stan.Msg) {
		log.Println("category event, msg info:", m)
		domain.GetConnection(utilities.GetConfiguration())
		m.Ack()
	}(msg)
}
