package enums

type WatcherType int

const (
	Watcher_HTTP WatcherType = iota
)

func (w WatcherType) String() string {
	switch w {
	case Watcher_HTTP:
		return "HTTP"
	default:
		return "Unknown"
	}
}
