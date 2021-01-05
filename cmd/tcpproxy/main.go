/*
 * Copyright (c) 2021 wuqingtao
 * tcpproxy is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 *
 * Author: wuqingtao (wqt_1110@qq.com)
 * CreateTime: 2021-01-05 21:41:39+0800
 */

package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tcpproxy/internal/proxy"
	"time"
)

var (
	addr     = flag.String("addr", ":8123", "代理的监听地址")
	backend  = flag.String("backend", "", "后端服务地址")
	certFile = flag.String("cert", "", "tls证书文件路径(pem)")
	keyFile  = flag.String("key", "", "tls密钥文件路径(pem)")
	debug    = flag.Bool("debug", false, "是否打印debug信息")
)

func main() {
	flag.Parse()
	if *addr == "" {
		log.Fatalln("代理的监听地址不能为空")
	}
	if *backend == "" {
		log.Fatalln("后端服务地址不能是为空")
	}

	// 初始化代理
	p := proxy.Proxy{
		Addr:    *addr,
		Backend: *backend,
		Dialer: &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 20 * time.Second,
		},
		Debug: *debug,
	}

	// 如果参数指定的tls证书和密钥，则读取证书并设置代理的TLSConfig
	if *certFile != "" && *keyFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("tls.LoadX509KeyPair(%s,%s) %s", *certFile, *keyFile, err)
		}
		p.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	// 开始监听
	go func() {
		if err := p.Listen(); err != nil {
			log.Fatalln(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, syscall.SIGHUP, syscall.SIGTERM)
	for s := range c {
		switch s {
		case os.Kill:
			log.Printf("接收到信号[Kill], 程序退成")
		default:
			log.Printf("接收到信号[%s]", s)
		}
	}
}
