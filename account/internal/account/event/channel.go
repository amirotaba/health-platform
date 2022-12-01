package event

import "git.paygear.ir/giftino/account/internal/account/entity"

const (
	ChannelCreatedEventType     = "create:channel"
	ChannelTypeUpdatedEventType = "updateType:channel"
	ChannelChargedEventType     = "channelCharged:channel"
	ChannelUpdatedEventType     = "update:channel"
	ChannelDeletedEventType     = "delete:channel"
)

type ChannelCreatedEvent struct {
	BaseEvent
	Data entity.CreateChannelRequest
}

type ChannelDepositedEvent struct {
	BaseEvent
}

type ChannelDeletedEvent struct {
	BaseEvent
	ChannelId string
}

type ChannelUpdatedEvent struct {
	BaseEvent
	Data entity.UpdateChannelRequest
}

type ChannelChargedEvent struct {
	BaseEvent
	Data entity.ChargeChannelRequest
}
