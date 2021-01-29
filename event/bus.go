package event

import "github.com/asaskevich/EventBus"

var busUtil = EventBus.New()

func Publish(topic string, args ...interface{}) {
	busUtil.Publish(topic, args...)
}
func SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return busUtil.SubscribeAsync(topic, fn, transactional)
}
func Subscribe(topic string, fn interface{}) error {
	return busUtil.Subscribe(topic, fn)
}
func Unsubscribe(topic string, fn interface{}) error {
	return busUtil.Unsubscribe(topic, fn)
}
func WaitAsync() {
	busUtil.WaitAsync()
}
