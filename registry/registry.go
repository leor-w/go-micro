// Package registry 为服务发现的接口
package registry

import (
	"errors"
)

var (
	DefaultRegistry = NewRegistry()

	// 调用 GetService 时没有找到对应服务时的错误
	ErrNotFound = errors.New("service not found")
	// 当 Watcher 停止时发生的错误
	ErrWatcherStopped = errors.New("watcher stopped")
)

// 注册中心为服务发现提供接口，并对不同的实现进行抽象
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

type DeregisterOption func(*DeregisterOptions)

type GetOption func(*GetOptions)

type ListOption func(*ListOptions)

// Register 注册一个服务节点。另外，可以提供诸如 TTL 之类的选项。
func Register(s *Service, opts ...RegisterOption) error {
	return DefaultRegistry.Register(s, opts...)
}

// Deregister 注销服务节点
func Deregister(s *Service) error {
	return DefaultRegistry.Deregister(s)
}

// GetService 检索服务。由于分割了 Name/Version，所以返回的是 Service 的切片。
func GetService(name string) ([]*Service, error) {
	return DefaultRegistry.GetService(name)
}

// 列出服务列表，只返回服务的名称
func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

// Watch 返回一个监视器，可以用来跟踪注册表的更新。
func Watch(opts ...WatchOption) (Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}

func String() string {
	return DefaultRegistry.String()
}
