package micro

// micro service 模块的实现

import (
	"os"
	"os/signal"
	rtime "runtime"
	"strings"
	"sync"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/debug/handler"
	"github.com/asim/go-micro/v3/debug/stats"
	"github.com/asim/go-micro/v3/debug/trace"
	"github.com/asim/go-micro/v3/logger"
	plugin "github.com/asim/go-micro/v3/plugins"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/store"
	signalutil "github.com/asim/go-micro/v3/util/signal"
	"github.com/asim/go-micro/v3/util/wrapper"
)

type service struct {
	opts Options

	once sync.Once
}

func newService(opts ...Option) Service {
	service := new(service)
	options := newOptions(opts...)

	// 服务名称
	serviceName := options.Server.Options().Name

	// 包装客户端以在任何调用中注入From-Service头
	options.Client = wrapper.FromService(serviceName, options.Client)
	options.Client = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)

	// 包装服务器以提供处理程序统计
	options.Server.Init(
		server.WrapHandler(wrapper.HandlerStats(stats.DefaultStats)),
		server.WrapHandler(wrapper.TraceHandler(trace.DefaultTracer)),
	)

	// 设置服务的选项
	service.opts = options

	return service
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init 初始化 options，此外，还会调用 cmd.init，它还会解析命令行的 flags。 cmd.Init 只会在首次执行 Init 时调用。
func (s *service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		// 设置 plugin
		for _, p := range strings.Split(os.Getenv("MICRO_PLUGIN"), ",") {
			if len(p) == 0 {
				continue
			}

			// 加载 plugin
			c, err := plugin.Load(p)
			if err != nil {
				logger.Fatal(err)
			}

			// 初始化 plugin
			if err := plugin.Init(c); err != nil {
				logger.Fatal(err)
			}
		}

		// 设置 cmd 名称
		if len(s.opts.Cmd.App().Name) == 0 {
			s.opts.Cmd.App().Name = s.Server().Options().Name
		}

		// 初始化命令的 flags, 覆盖新服务
		if err := s.opts.Cmd.Init(
			cmd.Auth(&s.opts.Auth),
			cmd.Broker(&s.opts.Broker),
			cmd.Registry(&s.opts.Registry),
			cmd.Runtime(&s.opts.Runtime),
			cmd.Transport(&s.opts.Transport),
			cmd.Client(&s.opts.Client),
			cmd.Config(&s.opts.Config),
			cmd.Server(&s.opts.Server),
			cmd.Store(&s.opts.Store),
			cmd.Profile(&s.opts.Profile),
		); err != nil {
			logger.Fatal(err)
		}

		// 显示的将表名设置为服务名
		name := s.opts.Cmd.App().Name
		s.opts.Store.Init(store.Table(name))
	})
}

// ########################################################## //
// #################### service 接口的实现 #################### //
// ########################################################## //

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) String() string {
	return "micro"
}

func (s *service) Start() error {
	// 执行启动前的钩子函数
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	// 服务启动
	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	// 执行启动后执行的钩子函数
	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	var gerr error

	// 执行停止前执行的钩子函数
	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	// 服务终止
	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	// 执行停止后执行的钩子函数
	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	// 注册调试处理
	s.opts.Server.Handle(
		s.opts.Server.NewHandler(
			handler.NewHandler(s.opts.Client),
			server.InternalHandler(true),
		),
	)

	// 启动服务分析器
	if s.opts.Profile != nil {
		// 查看互斥锁的争用情况
		rtime.SetMutexProfileFraction(5)
		// 查看阻塞配置文件
		rtime.SetBlockProfileRate(1)

		if err := s.opts.Profile.Start(); err != nil {
			return err
		}
		defer s.opts.Profile.Stop()
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("Starting [service] %s", s.Name())
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, signalutil.Shutdown()...)
	}

	select {
	// 等待关闭信号
	case <-ch:
	// 等待上下文退出信号
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

// ########################################################## //
// #################### service 接口的实现 #################### //
// ########################################################## //
