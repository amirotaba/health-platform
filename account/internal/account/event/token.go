package event

import "git.paygear.ir/giftino/account/internal/account/entity"

const (
	TokenCreatedEventType   = "create:token"
	TokenRefreshedEventType = "refresh:token"
)

type TokenCreatedEvent struct {
	BaseEvent
	Data entity.TokenEntity
}

type RefreshCreatedEvent struct {
	BaseEvent
	Data entity.UpdateTokenRequest
}
