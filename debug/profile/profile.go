// Package profile is for profilers
package profile

type Profile interface {
	// Start 启动性能分析
	Start() error
	// Stop 停止性能分析
	Stop() error
	// Name 分析器的名字
	String() string
}

var (
	DefaultProfile Profile = new(noop)
)

type noop struct{}

func (p *noop) Start() error {
	return nil
}

func (p *noop) Stop() error {
	return nil
}

func (p *noop) String() string {
	return "noop"
}

type Options struct {
	// Name to use for the profile
	Name string
}

type Option func(o *Options)

// Name of the profile
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
