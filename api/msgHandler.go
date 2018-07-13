package api

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/nats-io/go-nats-streaming"
)

const (
	kclusterID          = "api-cluster"
	kclientPublisherID  = "publisher-event-store"
	kclientSubscriberID = "subscriber-event-store"
)

//MessageHanderInfo configuration information
type MessageHanderInfo struct {
	clusterID string
	clientID  string
	msgURL    string
}

// DefaultMessagePublisherInfo represent default message for publisher
var defaultMessagePublisherInfo = MessageHanderInfo{
	clusterID: kclusterID,
	clientID:  kclientPublisherID,
	msgURL:    stan.DefaultNatsURL,
}

// DefaultMessageSubscriberInfo represent default message for publisher
var defaultMessageSubscriberInfo = MessageHanderInfo{
	clusterID: kclusterID,
	clientID:  kclientSubscriberID,
	msgURL:    stan.DefaultNatsURL,
}

//MessageHandler represent handler
type MessageHandler interface {
	Close() error
}

//MessagePublisher represent publiser message server
type MessagePublisher struct {
	msgConn *stan.Conn
}

//MessageSubscriber represent subscriber message server
type MessageSubscriber struct {
	msgConn          *stan.Conn
	lastProcessedSeq uint64
}

func createMessageConnection(defaultMsgInfo MessageHanderInfo, cfg MessageHanderInfo) (stan.Conn, error) {
	clusterID, clientID, msgURL := defaultMsgInfo.clusterID, defaultMsgInfo.clientID, defaultMsgInfo.msgURL
	if cfg.clientID != "" {
		clientID = cfg.clientID
	}
	if cfg.clusterID != "" {
		clusterID = cfg.clusterID
	}
	if cfg.msgURL != "" {
		msgURL = cfg.msgURL
	}
	return stan.Connect(clusterID, clientID, stan.NatsURL(msgURL))
}

// CreateMessagePublisher create MessagePublisher instance handler
func CreateMessagePublisher(cfg MessageHanderInfo) (*MessagePublisher, error) {
	sc, err := createMessageConnection(defaultMessagePublisherInfo, cfg)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	m := MessagePublisher{
		msgConn: &sc,
	}
	return &m, nil
}

// Close close resource
func (m *MessagePublisher) Close() error {
	conn := m.msgConn
	log.Println("closing msg publisher connection")
	if err := (*conn).Close(); err != nil {
		log.Println("error while closing msg server connection,", err)
		return err
	}
	return nil
}

// PublishEvent publish
func (m *MessagePublisher) PublishEvent(channel string, msg string) {
	publishMsgHandler := func(uid string, err error) {
		if err != nil {
			log.Printf("publishMsgHandler error msg, %v\n", err)
		} else {
			log.Printf("Received ACK for message id %s\n", uid)
		}
	}
	sc := *m.msgConn
	if uid, err := sc.PublishAsync(channel, []byte(msg), publishMsgHandler); err != nil {
		log.Printf("error publishing msg %v, error:%v, uid:%v", msg, err, uid)
	}
}

// CreateMessageSubscriber create MessageSubscriber instance handler
func CreateMessageSubscriber(cfg MessageHanderInfo) (*MessageSubscriber, error) {
	sc, err := createMessageConnection(defaultMessageSubscriberInfo, cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	m := MessageSubscriber{
		msgConn: &sc,
	}
	return &m, nil
}

//Close close resource
func (m *MessageSubscriber) Close() error {
	conn := m.msgConn
	log.Println("closing msg subscriber connection")
	if err := (*conn).Close(); err != nil {
		log.Println("error while closing msg server connection,", err)
		return err
	}
	return nil
}

//SubscribeEvent subscribe
func (m *MessageSubscriber) SubscribeEvent(channel string, durableID string, fn stan.MsgHandler) {
	handler := func(msg *stan.Msg) {
		log.Println("msg info:", msg)
		if m.lastProcessedSeq == 0 {
			log.Println("start up")
			//initially start and require logic to setup latest message sequence.
			atomic.SwapUint64(&m.lastProcessedSeq, msg.Sequence)
			return
		}
		if msg.Redelivered {
			log.Println("redelivered")
			atomic.SwapUint64(&m.lastProcessedSeq, msg.Sequence)
			//do process for redeliverd message

			msg.Ack()
			return
		}
		log.Println("new msg, procesing data")
		//do process for new message...
		msg.Ack()
		atomic.SwapUint64(&m.lastProcessedSeq, msg.Sequence)
	}
	//if durableID is nil, it will assign default handler above (example implementation).
	if fn != nil {
		handler = fn
	}
	sc := *m.msgConn
	wait, _ := time.ParseDuration("10s")
	sc.Subscribe(channel,
		handler,
		stan.DurableName(durableID),
		stan.MaxInflight(20),
		stan.SetManualAckMode(),
		stan.AckWait(wait))
}
