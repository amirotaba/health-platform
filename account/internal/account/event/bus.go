package event

type Bus interface {
	Publish(event []Event) error
	Register(topic string, handler func(event BaseEvent)) error
	ListTopic() []string
}
