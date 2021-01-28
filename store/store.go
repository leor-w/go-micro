// Package store is an interface for distributed data storage.
// The design document is located at https://github.com/micro/development/blob/master/design/store.md
package store

import (
	"errors"
	"time"
)

var (
	// ErrNotFound is returned when a key doesn't exist
	ErrNotFound = errors.New("not found")
	// DefaultStore is the memory store.
	DefaultStore Store = NewStore()
)

// Store 是一个数据存储接口
type Store interface {
	// Init 初始化存储。它必须在后备存储实现上执行任何必需的设置，并检查它是否可以使用，并返回任何错误。
	Init(...Option) error
	// Options 允许您查看当前选项。
	Options() Options
	// Read 接受单个键名和可选的 ReadOptions。它返回匹配的 []*Record 或错误。
	Read(key string, opts ...ReadOption) ([]*Record, error)
	// Write() 将记录写入存储，如果没有写入记录则返回错误。
	Write(r *Record, opts ...WriteOption) error
	// Delete 从存储中删除具有相应键值的记录。
	Delete(key string, opts ...DeleteOption) error
	// List 返回任何匹配的键，如果没有匹配则返回无错误的空列表。
	List(opts ...ListOption) ([]string, error)
	// Close 关闭存储。
	Close() error
	// String 返回实现的名称。
	String() string
}

// Record is an item stored or retrieved from a Store
type Record struct {
	// The key to store the record
	Key string `json:"key"`
	// The value within the record
	Value []byte `json:"value"`
	// Any associated metadata for indexing
	Metadata map[string]interface{} `json:"metadata"`
	// Time to expire a record: TODO: change to timestamp
	Expiry time.Duration `json:"expiry,omitempty"`
}

func NewStore(opts ...Option) Store {
	return NewMemoryStore(opts...)
}
