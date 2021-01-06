/*
 * Copyright (c) 2021 wuqingtao
 * tcpproxy is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 *
 * Author: wuqingtao (wqt_1110@qq.com)
 * CreateTime: 2021-01-05 21:41:39+0800
 */

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcpproxy/internal/proxy"
	"time"
)

var (
	addr      = flag.String("addr", ":8123", "代理的监听地址")
	backend   = flag.String("backend", "", "后端服务地址")
	certFile  = flag.String("tls.cert", "", "tls证书文件路径(pem格式)")
	keyFile   = flag.String("tls.key", "", "tls密钥文件路径(pem格式)")
	timeout   = flag.Int64("tcp.timeout", 5, "连接超时时间(单位秒)")
	keepalive = flag.Int64("tcp.keepalive", 15, "保持连接的时间间隔(单位秒)")
	debug     = flag.Bool("debug", false, "是否打印debug信息")
)

func main() {
	flag.Parse()
	if *addr == "" {
		log.Fatalln("代理的监听地址不能为空")
	}
	if *backend == "" {
		log.Fatalln("后端服务地址不能是为空")
	}

	p, err := proxy.NewProxy(*addr, *backend,
		proxy.WithTLSKeyPair(*certFile, *keyFile),
		proxy.WithTimeout(time.Duration(*timeout)*time.Second),
		proxy.WithKeepAlive(time.Duration(*keepalive)*time.Second),
		proxy.WithDebug(*debug),
	)
	if err != nil {
		log.Fatalf("proxy.NewProxy() error %s", err)
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
