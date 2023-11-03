package middleware

type EventBus interface {
}

type DefaultEventBus struct {
	eventStream chan struct{}
}
