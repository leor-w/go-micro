// Package runtime is a service runtime manager
package runtime

import (
	"errors"
	"time"
)

var (
	// DefaultRuntime is default micro runtime
	DefaultRuntime Runtime = NewRuntime()
	// DefaultName is default runtime service name
	DefaultName = "go.micro.runtime"

	ErrAlreadyExists = errors.New("already exists")
)

// Runtime 为服务运行时管理器
type Runtime interface {
	// Init 初始化 runtime。
	Init(...Option) error
	// Create 创建与注册服务。
	Create(*Service, ...CreateOption) error
	// Read 返回服务
	Read(...ReadOption) ([]*Service, error)
	// Update 在适当的时机更新服务。
	Update(*Service, ...UpdateOption) error
	// 删除一个服务
	Delete(*Service, ...DeleteOption) error
	// Logs 返回服务的日志
	Logs(*Service, ...LogsOption) (LogStream, error)
	// Start 启动 runtime
	Start() error
	// Stop 停止 runtime
	Stop() error
	// String runtime 的描述
	String() string
}

// Stream returns a log stream
type LogStream interface {
	Error() error
	Chan() chan LogRecord
	Stop() error
}

type LogRecord struct {
	Message  string
	Metadata map[string]string
}

// Scheduler is a runtime service scheduler
type Scheduler interface {
	// Notify publishes schedule events
	Notify() (<-chan Event, error)
	// Close stops the scheduler
	Close() error
}

// EventType defines schedule event
type EventType int

const (
	// Create is emitted when a new build has been craeted
	Create EventType = iota
	// Update is emitted when a new update become available
	Update
	// Delete is emitted when a build has been deleted
	Delete
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event is notification event
type Event struct {
	// ID of the event
	ID string
	// Type is event type
	Type EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service the event relates to
	Service *Service
	// Options to use when processing the event
	Options *CreateOptions
}

// Service is runtime service
type Service struct {
	// Name of the service
	Name string
	// Version of the service
	Version string
	// url location of source
	Source string
	// Metadata stores metadata
	Metadata map[string]string
}
