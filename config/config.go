// Package config is an interface for dynamic configuration.
package config

import (
	"context"

	"github.com/asim/go-micro/v3/config/loader"
	"github.com/asim/go-micro/v3/config/reader"
	"github.com/asim/go-micro/v3/config/source"
	"github.com/asim/go-micro/v3/config/source/file"
)

// Config 是动态配置的接口抽象
type Config interface {
	// 提供 reader.Values 的接口
	reader.Values
	// Init 初始化配置
	Init(opts ...Option) error
	// Options 配置中的选项
	Options() Options
	// 停止配置加载器/监视程序
	Close() error
	// Load 加载配置资源
	Load(source ...source.Source) error
	// 强制同步源变更集
	Sync() error
	// Watch 观察值的变化
	Watch(path ...string) (Watcher, error)
}

// Watcher is the config watcher
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}

type Option func(o *Options)

var (
	// Default Config Manager
	DefaultConfig, _ = NewConfig()
)

// NewConfig returns new config
func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

// Return config as raw json
func Bytes() []byte {
	return DefaultConfig.Bytes()
}

// Return config as a map
func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

// Force a source changeset sync
func Sync() error {
	return DefaultConfig.Sync()
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return DefaultConfig.Get(path...)
}

// Load config sources
func Load(source ...source.Source) error {
	return DefaultConfig.Load(source...)
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	return DefaultConfig.Watch(path...)
}

// LoadFile is short hand for creating a file source and loading it
func LoadFile(path string) error {
	return Load(file.NewSource(
		file.WithPath(path),
	))
}
