// Package micro 是用于微服务的可插入框架
package micro

import (
	"context"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/server"
)

type serviceKey struct{}

// Service 是一个接口，它将较为底层的库包装在 go-micro 中。这是一种构建和初始化服务比较方便的一种方法。
type Service interface {
	// service 的名称
	Name() string
	// Init 初始化 options
	Init(...Option)
	// Options 返回当前的 options
	Options() Options
	// Client 是用户调用服务的入口
	Client() client.Client
	// Server 是用来处理请求和事件的
	Server() server.Server
	// Run 运行服务
	Run() error
	// 服务的名称
	String() string
}

// Function 为单次执行 Service
type Function interface {
	// 继承 Service 的接口
	Service
	// Done 执行完成的信号
	Done() error
	// Handle 注册 RPC 处理程序
	Handle(v interface{}) error
	// Subscribe 注册一个订阅者
	Subscribe(topic string, v interface{}) error
}

/*
// Type Event is a future type for acting on asynchronous events
type Event interface {
	// Publish publishes a message to the event topic
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
	// Subscribe to the event
	Subscribe(ctx context.Context, v in
}

// Resource is a future type for defining dependencies
type Resource interface {
	// Name of the resource
	Name() string
	// Type of resource
	Type() string
	// Method of creation
	Create() error
}
*/

// Event 用于向主题发布消息
type Event interface {
	// Publish 向事件主题发布消息
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

// 类型别名以满足弃用要求
type Publisher = Event

type Option func(*Options)

var (
	HeaderPrefix = "Micro-"
)

// NewService 根据 option 创建并返回一个新的服务。
func NewService(opts ...Option) Service {
	return newService(opts...)
}

// FromContext 从上下文中检索服务。
func FromContext(ctx context.Context) (Service, bool) {
	s, ok := ctx.Value(serviceKey{}).(Service)
	return s, ok
}

// NewContext 返回嵌入服务的新上下文。
func NewContext(ctx context.Context, s Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, s)
}

// NewFunction 为一个单次执行的 Service 返会一个 Function
func NewFunction(opts ...Option) Function {
	return newFunction(opts...)
}

// NewEvent 创建一个新的事件发布者
func NewEvent(topic string, c client.Client) Event {
	if c == nil {
		c = client.NewClient()
	}
	return &event{c, topic}
}

// 已弃用: NewPublisher 返回一个新的 Publisher
func NewPublisher(topic string, c client.Client) Event {
	return NewEvent(topic, c)
}

// RegisterHandler 是用于注册处理程序的语法糖
func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}

// RegisterSubscriber 是用于注册订户的语法糖
func RegisterSubscriber(topic string, s server.Server, h interface{}, opts ...server.SubscriberOption) error {
	return s.Subscribe(s.NewSubscriber(topic, h, opts...))
}
