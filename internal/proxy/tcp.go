/*
 * Copyright (c) 2021 wuqingtao
 * proxy is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
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
