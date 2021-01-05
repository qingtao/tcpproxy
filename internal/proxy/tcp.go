// koroFileHeader generate.
// Author: 吴庆涛
// Description: tcp proxy
// Date: 2021-01-05 19:37:22
// FilePath: /tcpproxy/internal/proxy/tcp.go
// LastEditors: 吴庆涛
// LastEditTime: 2021-01-05 20:45:18

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
	Addr string
	// 后端地址
	Backend string
	// tls配置
	TLSConfig *tls.Config
	// 连接器
	Dialer *net.Dialer
	// 是否打印调试
	Debug bool
}

// Listen 监听tcp端口
func (p *Proxy) Listen() error {
	l, err := net.Listen("tcp", p.Addr)
	if err != nil {
		return err
	}
	if p.TLSConfig != nil {
		l = tls.NewListener(l, p.TLSConfig)
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
		if p.Debug {
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
	agentConn, err := p.Dialer.DialContext(ctx, "tcp", p.Backend)
	if err != nil {
		handleError(fmt.Errorf("client %s connect failed %w", clientConn.RemoteAddr(), err))
		return
	}

	defer func() {
		agentConn.Close()
		if p.Debug {
			log.Printf("agent client %s <-> server %s closed",
				agentConn.LocalAddr(),
				agentConn.RemoteAddr(),
			)
		}
	}()

	if p.Debug {
		log.Printf("client %s <-> agent server %s - agent client %s <-> server %s connected",
			clientConn.RemoteAddr(),
			clientConn.LocalAddr(),
			agentConn.LocalAddr(),
			agentConn.RemoteAddr(),
		)
	}
	copyData(clientConn, agentConn)
}
