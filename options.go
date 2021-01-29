package micro

import (
	"context"
	"time"

	"github.com/asim/go-micro/v3/auth"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/config"
	"github.com/asim/go-micro/v3/debug/profile"
	"github.com/asim/go-micro/v3/debug/trace"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/runtime"
	"github.com/asim/go-micro/v3/selector"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/store"
	"github.com/asim/go-micro/v3/transport"
	"github.com/micro/cli/v2"
)

// Options 微服务的选项
type Options struct {
	// 权限选项
	Auth auth.Auth
	// 异步消息传递选项
	Broker broker.Broker
	// 解析命令的接口
	Cmd cmd.Cmd
	// 配置中心接口
	Config config.Config
	// 客户端，用于向服务发出请求的接口
	Client client.Client
	// 微服务的抽象
	Server server.Server
	// 存储的抽象
	Store store.Store
	// 服务注册的接口，用于服务发现
	Registry registry.Registry
	// 服务运行时管理器。
	Runtime runtime.Runtime
	// 用于服务之间通信的接口。
	Transport transport.Transport
	// 性能分析器
	Profile profile.Profile

	// 在启动、停止之前以及启动、停止之后执行的钩子
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// 接口实现的其他选项可以存储在上下文中
	Context context.Context

	Signal bool
}

// newOptions 创建微服务的选项
func newOptions(opts ...Option) Options {
	opt := Options{
		Auth:      auth.DefaultAuth,
		Broker:    broker.DefaultBroker,
		Cmd:       cmd.DefaultCmd,
		Config:    config.DefaultConfig,
		Client:    client.DefaultClient,
		Server:    server.DefaultServer,
		Store:     store.DefaultStore,
		Registry:  registry.DefaultRegistry,
		Runtime:   runtime.DefaultRuntime,
		Transport: transport.DefaultTransport,
		Context:   context.Background(),
		Signal:    true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Broker 用于微服务代理
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
		// 更新 客户端 和 服务器
		o.Client.Init(client.Broker(b))
		o.Server.Init(server.Broker(b))
	}
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

// Client 用于服务的客户端
func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Context 为服务指定上下文。可以用来表示服务的关闭，并获得额外的选项值。
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// HandleSignal 切换信号处理程序的自动安装，改处理程序捕捉 TERM、INT 与 QUIT。
// 使用此功能禁用信号处理程序的用户，应通过 context 控制服务的活跃度
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

// Profile 设置用于调试用的 profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.Profile = p
	}
}

// Server 设置服务的 option
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Store 设置 store 的 option
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// Registry 为服务和底层组件设置注册中心
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		// 更新 客户端 和 服务器
		o.Client.Init(client.Registry(r))
		o.Server.Init(server.Registry(r))
		// Update Broker

		// 更新代理
		o.Broker.Init(broker.Registry(r))
	}
}

// Tracer 设置服务的 跟踪程序
func Tracer(t trace.Tracer) Option {
	return func(o *Options) {
		o.Server.Init(server.Tracer(t))
	}
}

// Auth 设置服务的身份验证器
func Auth(a auth.Auth) Option {
	return func(o *Options) {
		o.Auth = a
	}
}

// Config 设置服务的 配置中心
func Config(c config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// Selector 设置服务 客户端的 选择器
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Client.Init(client.Selector(s))
	}
}

// Transport 设置服务和底层组件的传输
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		// 更新到服务器与客户端
		o.Client.Init(client.Transport(t))
		o.Server.Init(server.Transport(t))
	}
}

// Runtime 设置运行时管理器
func Runtime(r runtime.Runtime) Option {
	return func(o *Options) {
		o.Runtime = r
	}
}

// 便利选项设置

// Address 设置服务器地址
func Address(addr string) Option {
	return func(o *Options) {
		o.Server.Init(server.Address(addr))
	}
}

// Name 设置服务的名称
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}

// Version 设置服务的版本
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Init(server.Version(v))
	}
}

// Metadata 将元数据与服务关联
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Server.Init(server.Metadata(md))
	}
}

// Flags 将 flags 传递给服务
func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

// Action 操作可以用于解析用户提供的 cli 选项
func Action(a func(*cli.Context) error) Option {
	return func(o *Options) {
		o.Cmd.App().Action = a
	}
}

// RegisterTTL 注册服务是要使用的 TTL 值
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.Server.Init(server.RegisterTTL(t))
	}
}

// RegisterInterval 指定服务重新注册的时间间隔
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.Server.Init(server.RegisterInterval(t))
	}
}

// WrapClient 是一种用中间件包装 Client 的方式。可以提供包装器的列表。包装器器是按照先进后出方式执行的，因此最后一个包装器是最后执行的。
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall 包装了 Client 调用该函数
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler 将处理函数包装器添加并传递到服务器的选项列表。
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// 单次初始化
		o.Server.Init(wrappers...)
	}
}

// WrapSubscriber 将订阅者包装器添加并传递到服务器的选项列表中。
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// 单次初始化
		o.Server.Init(wrappers...)
	}
}

// 执行前与执行后

// BeforeStart 执行服务启动前的钩子函数
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop 执行服务停止前的钩子函数
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart 执行服务启动后的钩子函数
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop 执行服务停止后的钩子函数
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
