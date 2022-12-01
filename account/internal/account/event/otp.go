package event

import "git.paygear.ir/giftino/account/internal/account/entity"

const (
	OtpCreatedEventType       = "create:otp"
	OtpStatusUpdatedEventType = "updateStatus:otp"
	OtpDeletedEventType       = "delete:otp"
)

type OtpCreatedEvent struct {
	BaseEvent
	entity.OtpEntity
}
