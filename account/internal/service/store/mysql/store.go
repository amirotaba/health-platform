package mysql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/event"
)

type store struct {
	db        *sql.DB
	EventChan chan event.Event
}

func NewStore(db *sql.DB) event.Store {
	return &store{
		db:        db,
		EventChan: make(chan event.Event),
	}
}

func (s *store) Store(events []event.Event) error {
	for _, e := range events {
		err := s.saveEvent(e)
		if err != nil {
			return err
		}

		// todo handle this part must be complete
		s.EventChan <- e
	}

	return nil
}

func (s *store) Load(uuid string) ([]event.Event, error) {
	rows, err := s.db.Query("SELECT * FROM events WHERE aggregate_id=? ORDER BY id ASC;", uuid)
	// declare empty accounts variable
	if err != nil {
		return []event.Event{}, err
	}

	defer rows.Close()
	var events []event.Event
	// iterate over rows
	for rows.Next() {
		var e event.BaseEventModel
		err = rows.Scan(
			&e.ID,
			&e.AggregateID,
			&e.AggregateType,
			&e.Type,
			&e.Version,
			&e.Tag,
			&e.Data,
			&e.CreatedAt,
		)

		log.Println(err)
		if err != nil {
			return events, err
		}

		events = append(events, newEvent(e))
	}

	return events, nil

}

func (s *store) saveEvent(e event.Event) error {
	stmt := `insert into events(
								aggregate_id,
	                            aggregate_type,
                                type,
	                            version,
                                tag,
	                            data) values(?,?,?,?,?,?);`
	rows, err := s.db.Query(stmt,
		e.GetAggregateID(),
		e.GetAggregateType(),
		e.GetEventType(),
		e.GetVersion(),
		e.GetTag(),
		e.GetEventData())
	if err != nil {
		return err
	}

	defer rows.Close()
	var id sql.NullInt64
	for rows.Next() {
		err = rows.Scan(
			&id,
		)
		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) GetEventChan() <-chan event.Event {
	return s.EventChan
}

func newEvent(model event.BaseEventModel) event.BaseEvent {
	return event.BaseEvent{
		ID:            model.ID.Int64,
		AggregateID:   model.AggregateID.String,
		AggregateType: model.AggregateType.String,
		Type:          model.Type.String,
		Version:       model.Version.Int64,
		Data:          model.Data.String,
	}
}
