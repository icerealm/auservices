package api

import (
	"log"
	"os"
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

//GetEventSubscribers get eventsubscribers instance.
func GetEventSubscribers() *EventSubscriber {
	subMap := make(map[string]*MessageSubscriber)
	subscribers := []*MessageSubscriber{}
	//initial category subscriber
	categorySub, err := CreateMessageSubscriber(
		MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial category subscriber")
		os.Exit(1)
	}
	categorySub.SubscribeEvent(kcategoryChannelID, durableCategoryID, nil) //test
	subMap[durableCategoryID] = categorySub
	subscribers = append(subscribers, categorySub)

	eventSubscriber = &EventSubscriber{
		subscriberMap: subMap,
		Subscribers:   subscribers,
	}
	return eventSubscriber
}
