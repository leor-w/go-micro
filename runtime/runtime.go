// Package runtime 是一个运行服务时管理器
package runtime

import (
	"errors"
	"time"
)

var (
	// DefaultRuntime 默认的 micro 运行时服务
	DefaultRuntime Runtime = NewRuntime()
	// DefaultName 默认的运行时服务名称
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

// LogStream 返回一个日志流
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
// Scheduler 为运行时的服务调度程序
type Scheduler interface {
	// Notify 发布调度事件
	Notify() (<-chan Event, error)
	// Close 关闭服务调度
	Close() error
}

// EventType 定义调度事件类型
type EventType int

const (
	// Create 创建新构建时的类型
	Create EventType = iota
	// Update is emitted when a new update become available
	// 当新的更新可用时，发出的类型
	Update
	// Delete 删除构建时，发出的类型
	Delete
)

// String 返回事件类型名称
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

// Event 是通知事件
type Event struct {
	// ID 事件的ID
	ID string
	// Type 事件的类型
	Type EventType
	// Timestamp 事件的事件戳
	Timestamp time.Time
	// Service 与事件关联的服务
	Service *Service
	// Options 处理事件的选项
	Options *CreateOptions
}

// Service 运行时服务
type Service struct {
	// Name 服务名称
	Name string
	// Version 服务版本
	Version string
	// URL 源的位置
	Source string
	// Metadata 元数据的储存对象
	Metadata map[string]string
}
