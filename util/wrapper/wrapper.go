package wrapper

import (
	"context"
	"strings"

	"github.com/asim/go-micro/v3/auth"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/debug/stats"
	"github.com/asim/go-micro/v3/debug/trace"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/server"
)

type fromServiceWrapper struct {
	client.Client

	// headers to inject
	headers metadata.Metadata
}

var (
	HeaderPrefix = "Micro-"
)

func (f *fromServiceWrapper) setHeaders(ctx context.Context) context.Context {
	// 不要覆盖键
	return metadata.MergeContext(ctx, f.headers, false)
}

func (f *fromServiceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Call(ctx, req, rsp, opts...)
}

func (f *fromServiceWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	ctx = f.setHeaders(ctx)
	return f.Client.Stream(ctx, req, opts...)
}

func (f *fromServiceWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Publish(ctx, p, opts...)
}

// FromService 包装客户端以注入服务和认证元数据
func FromService(name string, c client.Client) client.Client {
	return &fromServiceWrapper{
		c,
		metadata.Metadata{
			HeaderPrefix + "From-Service": name,
		},
	}
}

// HandlerStats 包装一个服务器处理程序来生成请求/错误统计信息
func HandlerStats(stats stats.Stats) server.HandlerWrapper {
	// 返回一个请求处理的包装
	return func(h server.HandlerFunc) server.HandlerFunc {
		// 返回一个返回函数的函数
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// 执行请求处理函数
			err := h(ctx, req, rsp)
			// 记录数据
			stats.Record(err)
			// 返回错误
			return err
		}
	}
}

type traceWrapper struct {
	client.Client

	name  string
	trace trace.Tracer
}

func (c *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	newCtx, s := c.trace.Start(ctx, req.Service()+"."+req.Endpoint())

	s.Type = trace.SpanTypeRequestOutbound
	err := c.Client.Call(newCtx, req, rsp, opts...)
	if err != nil {
		s.Metadata["error"] = err.Error()
	}

	// 结束服务跟踪
	c.trace.Finish(s)

	return err
}

// TraceCall 调用服务跟踪的包装
func TraceCall(name string, t trace.Tracer, c client.Client) client.Client {
	return &traceWrapper{
		name:   name,
		trace:  t,
		Client: c,
	}
}

// TraceHandler 包装服务器处理程序以执行跟踪
func TraceHandler(t trace.Tracer) server.HandlerWrapper {
	// 返回一个处理函数的包装
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// debug 不保存跟踪信息
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// get the span
			newCtx, s := t.Start(ctx, req.Service()+"."+req.Endpoint())
			s.Type = trace.SpanTypeRequestInbound

			err := h(newCtx, req, rsp)
			if err != nil {
				s.Metadata["error"] = err.Error()
			}

			// finish
			t.Finish(s)

			return err
		}
	}
}

// 权限包装
type authWrapper struct {
	client.Client
	auth func() auth.Auth
}

func (a *authWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// 解析选项
	var options client.CallOptions
	for _, o := range opts {
		o(&options)
	}

	// 检查是否已经设置了授权头。除非指定了 ServiceToken 选项或者没有提供该头，否则不会覆盖该头
	if _, ok := metadata.Get(ctx, "Authorization"); ok && !options.ServiceToken {
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// 如果 auth 为 nil，将无法获取访问令牌，所以我们在没有令牌的情况下执行请求。
	aa := a.auth()
	if aa == nil {
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// 如果没有设置命名空间则设置一个命名空间（例如 在服务请求中设置）
	if _, ok := metadata.Get(ctx, "Micro-Namespace"); !ok {
		ctx = metadata.Set(ctx, "Micro-Namespace", aa.Options().Namespace)
	}

	// 检查是否有一个有效的访问令牌
	aaOpts := aa.Options()
	if aaOpts.Token != nil && !aaOpts.Token.Expired() {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+aaOpts.Token.AccessToken)
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// 在没有身份验证令牌的情况下调用
	return a.Client.Call(ctx, req, rsp, opts...)
}
