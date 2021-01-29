package server

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/codec"
	"github.com/asim/go-micro/v3/debug/trace"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/transport"
)

type Options struct {
	Codecs       map[string]codec.NewCodec
	Broker       broker.Broker
	Registry     registry.Registry
	Tracer       trace.Tracer
	Transport    transport.Transport
	Metadata     map[string]string
	Name         string
	Address      string
	Advertise    string
	Id           string
	Version      string
	HdlrWrappers []HandlerWrapper
	SubWrappers  []SubscriberWrapper

	// RegisterCheck 注册服务前运行检查功能
	RegisterCheck func(context.Context) error
	// 寄存器失效时间
	RegisterTTL time.Duration
	// 注册的间隔事件
	RegisterInterval time.Duration

	// 请求的路由器
	Router Router

	// TLSConfig TLS 配置。
	TLSConfig *tls.Config

	// 接口实现的其他选项可以存储在上下文中
	Context context.Context
}

func newOptions(opt ...Option) Options {
	opts := Options{
		Codecs:           make(map[string]codec.NewCodec),
		Metadata:         map[string]string{},
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Broker == nil {
		opts.Broker = broker.DefaultBroker
	}

	if opts.Registry == nil {
		opts.Registry = registry.DefaultRegistry
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = DefaultRegisterCheck
	}

	if len(opts.Address) == 0 {
		opts.Address = DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = DefaultName
	}

	if len(opts.Id) == 0 {
		opts.Id = DefaultId
	}

	if len(opts.Version) == 0 {
		opts.Version = DefaultVersion
	}

	return opts
}

// Server 服务器名称
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// 唯一的服务器 ID
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Version 服务的版本
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address 服务的 地址与端口 host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// 用于服务发现的地址 - host:port
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Broker 用于发布和订阅的代理器
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

// Codec 用于对给定内容类型的请求进行编码/解码
func Codec(contentType string, c codec.NewCodec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// Context 为 service 指定一个上下文。用于给服务发送关闭信号，和为服务添加额外的选项值。
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Registry 用于服务发现
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Tracer 分布式的服务追踪
func Tracer(t trace.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Transport 通信机制，如 http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Metadata 关联元数据与服务器
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// RegisterCheck 在注册之前执行 fn。
func RegisterCheck(fn func(context.Context) error) Option {
	return func(o *Options) {
		o.RegisterCheck = fn
	}
}

// 指定注册服务的 TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// 注册服务的间隔
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// TLSConfig 指定一个 *tls.Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		// 设置到 options 内部的 tls
		o.TLSConfig = t

		// 如果 options 没有设置默认传输方式，则设置一个传输方式
		// 下面的 Init 需要设置一个默认的传输方式
		if o.Transport == nil {
			o.Transport = transport.DefaultTransport
		}

		// 设置 tls 到传输方式
		o.Transport.Init(
			transport.Secure(true),
			transport.TLSConfig(t),
		)
	}
}

// WithRouter 设置请求的路由
func WithRouter(r Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Wait 告诉服务器在退出之前要等待请求完成，如果 `wg` 为空，
// 服务器只等待 rpc 处理完成，如果需要更细粒度控制的用户，
// 需要在这里传递一个具体的 `wg`，服务器将在停止时等待它的完成。
func Wait(wg *sync.WaitGroup) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		if wg == nil {
			wg = new(sync.WaitGroup)
		}
		o.Context = context.WithValue(o.Context, "wait", wg)
	}
}

// WrapHandler 添加一个 handler 的包装器并传递到服务的 options 列表
func WrapHandler(w HandlerWrapper) Option {
	return func(o *Options) {
		o.HdlrWrappers = append(o.HdlrWrappers, w)
	}
}

// Adds a subscriber Wrapper to a list of options passed into the server
// WrapSubscriber 添加一个 subscriber 的包装器并传递到服务器的 options 列表
func WrapSubscriber(w SubscriberWrapper) Option {
	return func(o *Options) {
		o.SubWrappers = append(o.SubWrappers, w)
	}
}
