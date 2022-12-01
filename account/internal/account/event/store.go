package event

type Store interface {
	Store([]Event) error
	Load(uuid string) ([]Event, error)
	GetEventChan() <-chan Event
}
