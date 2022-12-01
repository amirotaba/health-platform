package event

type Handler interface {
	Register(topic string) error
	//HandleEvents()
}
