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
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
)

// Proxy 代理结构
type Proxy struct {
	// 服务监听地址
	Addr string
	// 后端地址
	Backend string
	// tls配置
	TLSConfig *tls.Config
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
	log.Printf("error %s\n", err)
}

func copyData(c1, c2 net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			handleError(fmt.Errorf("%w", e))
		}
	}()
	defer func() {
		if e := c1.Close(); e != nil {
			handleError(fmt.Errorf("%w", e))
		}
	}()
	defer func() {
		if e := c2.Close(); e != nil {
			handleError(fmt.Errorf("%w", e))
		}
	}()
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}

func (p *Proxy) handleConnect(conn net.Conn) {
	c, err := net.Dial("tcp", p.Backend)
	if err != nil {
		handleError(fmt.Errorf("client %s connect failed %w", conn.RemoteAddr(), err))
		return
	}
	var s string
	defer func() {
		log.Printf("%s closed", s)
	}()
	s = fmt.Sprintf("client %s <- %s - %s -> server %s",
		conn.RemoteAddr(),
		conn.LocalAddr(),
		c.LocalAddr(),
		c.RemoteAddr(),
	)
	log.Printf("%s connected", s)
	copyData(conn, c)
}
