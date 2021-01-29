package server

import (
	"context"
)

// HandlerFunc 请求处理的包装器。接收的参数为具体 request 与 response。
type HandlerFunc func(ctx context.Context, req Request, rsp interface{}) error

// SubscriberFunc 订阅服务器的包装器，接收的实际参数为具体发布的消息
type SubscriberFunc func(ctx context.Context, msg Message) error

// HandlerWrapper 包装一个 HandlerFunc 和返回一个对等的 HandlerFunc
type HandlerWrapper func(HandlerFunc) HandlerFunc

// SubscriberWrapper 包装一个 SubscriberFunc 和返回一个对等的 SubscriberFunc
type SubscriberWrapper func(SubscriberFunc) SubscriberFunc

// StreamWrapper 包装一个 Stream 接口并返回一个对等的 Stream。
// 因为流在方法调用的生命周期内都存在，可以包装服务跟踪、监控、统计等各种流
type StreamWrapper func(Stream) Stream
