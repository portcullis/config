package config

// Notifier for configuration Setting changes
type Notifier interface {
	// Notify defines a function that is called when s.Set is called with a different value other than the current
	Notify(s *Setting)
}

// NotifyHandle is used to stop notifications of Setting changes
type NotifyHandle struct {
	stopFunc func(interface{})
}

// Close the notification handle
func (h *NotifyHandle) Close() error {
	if h.stopFunc == nil {
		return nil
	}

	h.stopFunc(h)

	return nil
}

// NotifyFunc defines a function that is called when s.Set is called with a different value other than the current
type NotifyFunc func(s *Setting)

// Notify implements Notifier.Notify
func (f NotifyFunc) Notify(s *Setting) {
	f(s)
}
