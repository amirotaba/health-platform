package event_handler

//import (
//	"encoding/json"
//	"git.paygear.ir/giftino/account/internal/account/command"
//	"log"
//
//	"github.com/nats-io/nats.go"
//
//	"git.paygear.ir/giftino/account/internal/account/entity"
//	"git.paygear.ir/giftino/account/internal/account/event"
//)
//
//// EventHandler Event Handler subscribes 'persisted Events channel' (provided by event Store)
//// and reacts with commands for some of them
//type EventHandler struct {
//	store event.Store
//	conn  *nats.Conn
//}
//
//// NewEventHandler Constructor for new Event handler.
//func NewEventHandler(store event.Store, conn *nats.Conn) event.Handler {
//	return &EventHandler{store: store, conn: conn}
//}
//
//// Handle event logic
//func (eh *EventHandler) handleEvent(evt event.Event) {
//	switch evt.GetEventType() {
//	case event.ChannelCreatedEventType:
//		var payload entity.ChargeChannelRequest
//		err := json.Unmarshal([]byte(evt.GetEventData()), &payload)
//		if err != nil {
//			log.Println(err)
//		}
//
//		cmd := command.NewChargeChannelCommand(evt.GetAggregateID(), payload)
//
//	}
//}
//
//func (eh *EventHandler) Register(topic string) error {
//	if _, err := eh.conn.Subscribe(topic, func(m *nats.Msg) {
//		var baseEvent event.BaseEvent
//		err := json.Unmarshal(m.Data, &baseEvent)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		eh.handleEvent(baseEvent)
//	}); err != nil {
//		return err
//	}
//
//	return nil
//}
