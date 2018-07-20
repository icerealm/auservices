package msghandler

import (
	"auservices/api"
	"auservices/domain"
	"auservices/utilities"
	"encoding/json"
	"log"

	"github.com/lib/pq"
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

	//initial itemline subscriber
	itemlineSub := initItemLineSubscriber()
	eventSubscriber.AppendSubscriber(kitemLineChannelID, itemlineSub)

	return eventSubscriber
}

func initCategorySubscriber() *MessageSubscriber {
	//initial category subscriber
	sub, err := CreateMessageSubscriber(
		MessageHanderInfo{
			clientID: "category-sub-client",
		},
	)
	if err != nil {
		log.Fatalf("could not initial category subscriber")
	}
	sub.SubscribeEvent(kcategoryChannelID, durableCategoryID, categoryEventMsgHandler)
	return sub
}

func initItemLineSubscriber() *MessageSubscriber {
	//initial itemLine subscriber
	sub, err := CreateMessageSubscriber(
		MessageHanderInfo{
			clientID: "itemline-sub-client",
		},
	)
	if err != nil {
		log.Fatalf("could not initial itemLine subscriber")
	}
	sub.SubscribeEvent(kitemLineChannelID, durableCategoryID, nil)
	return sub
}

func shouldRejectError(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		log.Println("reject info, error:,", err, ", code:", err.Code)
		if err.Code == "23505" {
			return true
		}
	}
	return false
}

//categoryEventMsgHandler to handle business logic for category event
func categoryEventMsgHandler(msg *stan.Msg) {
	go func(m *stan.Msg) {
		db, err := domain.GetConnection(utilities.GetConfiguration())

		if err != nil {
			log.Fatalf("Category Event, error while connecting to database: %v", err)
		}
		defer db.Close()

		var category api.Category
		if err := json.Unmarshal(msg.Data, &category); err != nil {
			log.Printf("Category Event, error while converting message data: %v, msgInfo:%+v \n", err, msg)
			m.Ack()
		} else {
			err = domain.CreateCategory(db, msg.Sequence, &category, "cat-subscriber")
			if shouldRejectError(err) {
				log.Printf("Category Event, reject error and force ackknowledge msg sequence:%v \n", msg.Sequence)
			} else if err != nil {
				log.Printf("Category Event, error while inserting data: %v, msgInfo:%+v \n", err, msg)
			}
			m.Ack()
		}
	}(msg)
}

func itemLineEventMsgHandler(msg *stan.Msg) {
	go func(m *stan.Msg) {
		db, err := domain.GetConnection(utilities.GetConfiguration())
		if err != nil {
			log.Fatalf("ItemLine Event, error while connecting to database: %v", err)
		}
		defer db.Close()

	}(msg)
}
