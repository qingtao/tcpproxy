package proxy

import "time"

// Options 连接参数
type Options struct {
	// 服务监听地址
	addr string
	// 后端地址
	backend string
	// tls 证书路径
	cert string
	// tls 秘钥路径
	key string
	// 连接超时时间
	timeout time.Duration
	// 长连接保持时间
	keepAlive time.Duration
	// 是否打印调试
	debug bool
}

// Option 连接参数选项构造函数
type Option func(*Options)

// NewOptions 新建Options
func NewOptions(addr, backend string, options ...Option) *Options {
	opts := &Options{addr: addr, backend: backend}
	mergeOptions(opts, options...)
	return opts
}

func mergeOptions(opts *Options, options ...Option) {
	if opts == nil {
		panic("opts is nil")
	}
	for _, opt := range options {
		opt(opts)
	}
}

// WithTLSKeyPair 配置tls证书和秘钥的路径
func WithTLSKeyPair(certFile, keyFile string) Option {
	return func(o *Options) {
		o.cert = certFile
		o.key = keyFile
	}
}

// WithTimeout 配置超时参数
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.timeout = d
	}
}

// WithKeepAlive 配置长连接的时间间隔
func WithKeepAlive(d time.Duration) Option {
	return func(o *Options) {
		o.keepAlive = d
	}
}

// WithDebug 配置调试选项
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.debug = debug
	}
}
