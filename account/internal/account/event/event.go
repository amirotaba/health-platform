package event

import "database/sql"

const (
	Tag      = "event"
	SnapShot = "snapShot"
)

type PropertyMap map[string]interface{}

type Event interface {
	GetAggregateID() string
	GetAggregateType() string
	GetVersion() int64
	GetEventType() string
	GetEventData() string
	GetTag() string
}

// BaseEvent will save in mysql database and its data contains all info that cause the change
type BaseEvent struct {
	ID            int64  `json:"id"`
	AggregateID   string `json:"aggregate_id"`
	AggregateType string `json:"aggregate_type"`
	Type          string `json:"type"`
	Version       int64  `json:"version"`
	Tag           string `json:"tags"`
	Data          string `json:"data"`

	//Data          PropertyMap `json:"-"` // should save as jsonb
}

type BaseEventModel struct {
	ID            sql.NullInt64
	AggregateID   sql.NullString
	AggregateType sql.NullString // account, employee, affiliate
	Type          sql.NullString // create, delete, update
	Version       sql.NullInt64
	Tag           sql.NullString
	Data          sql.NullString
	CreatedAt     sql.NullTime
}

func (c BaseEvent) GetAggregateID() string {
	return c.AggregateID
}

func (c BaseEvent) GetAggregateType() string {
	return c.AggregateType
}

func (c BaseEvent) GetVersion() int64 {
	return c.Version
}

func (c BaseEvent) GetEventType() string {
	return c.Type
}

func (c BaseEvent) GetTag() string {
	return c.Tag
}

func (c BaseEvent) GetEventData() string {
	return c.Data
}
