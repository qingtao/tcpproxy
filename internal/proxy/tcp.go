/*
 * Copyright (c) 2021 wuqingtao
 * proxy is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 *
 * Author: wuqingtao (wqt_1110@qq.com)
 * CreateTime: 2021-01-05 21:41:39+0800
 */

// Package proxy 实现简单的tcp代理
package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
)

// Proxy 代理结构
type Proxy struct {
	// 服务监听地址
	addr string
	// 后端地址
	backend string
	// tls配置
	tlsConfig *tls.Config
	// 连接器
	dialer *net.Dialer
	// 是否打印调试
	debug bool
}

// NewProxy 新的代理
func NewProxy(addr, backend string, options ...Option) (*Proxy, error) {
	opts := NewOptions(addr, backend, options...)
	p := &Proxy{
		addr:    opts.addr,
		backend: opts.backend,
		debug:   opts.debug,
	}
	if opts.cert != "" && opts.key != "" {
		cert, err := tls.LoadX509KeyPair(opts.cert, opts.key)
		if err != nil {
			return nil, fmt.Errorf("tls.LoadX509KeyPair(%s,%s) %w", opts.cert, opts.key, err)
		}
		p.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	d := &net.Dialer{}
	if opts.timeout > 0 {
		d.Timeout = opts.timeout
	}
	if opts.keepAlive > 0 {
		d.KeepAlive = opts.keepAlive
	}
	p.dialer = d
	return p, nil
}

// Listen 监听tcp端口
func (p *Proxy) Listen() error {
	l, err := net.Listen("tcp", p.addr)
	if err != nil {
		return err
	}
	if p.tlsConfig != nil {
		l = tls.NewListener(l, p.tlsConfig)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			handleError(err)
		}
		go p.handleConnect(conn)
	}
}

func handleError(err error) {
	log.Printf("error %s", err)
}

func copyData(c1, c2 net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			handleError(fmt.Errorf("panic %w", e))
			debug.PrintStack()
		}
	}()
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}

func (p *Proxy) handleConnect(clientConn net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			handleError(fmt.Errorf("panic %w", e))
			debug.PrintStack()
		}
	}()
	defer func() {
		clientConn.Close()
		if p.debug {
			log.Printf("client %s <-> agent server %s closed",
				clientConn.RemoteAddr(),
				clientConn.LocalAddr(),
			)
		}
	}()
	if c, ok := clientConn.(*tls.Conn); ok {
		if err := c.Handshake(); err != nil {
			log.Printf("tls.Conn.Handshake() error %s", err)
			return
		}
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	agentConn, err := p.dialer.DialContext(ctx, "tcp", p.backend)
	if err != nil {
		handleError(fmt.Errorf("client %s connect failed %w", clientConn.RemoteAddr(), err))
		return
	}

	defer func() {
		agentConn.Close()
		if p.debug {
			log.Printf("agent client %s <-> server %s closed",
				agentConn.LocalAddr(),
				agentConn.RemoteAddr(),
			)
		}
	}()

	if p.debug {
		log.Printf("client %s <-> agent server %s - agent client %s <-> server %s connected",
			clientConn.RemoteAddr(),
			clientConn.LocalAddr(),
			agentConn.LocalAddr(),
			agentConn.RemoteAddr(),
		)
	}
	copyData(clientConn, agentConn)
}
