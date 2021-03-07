package pubsub

import "time"

// SystemUpDownEvent 系统启停事件
type SystemUpDownEvent struct {
	Up        bool
	CreatedAt time.Time
}