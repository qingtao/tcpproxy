# 一个简单的 tcp 代理

> 已知问题: 当前不支持后端服务是 tls 连接

```sh
Usage of tcpproxy:
  -addr string
        代理的监听地址 (default ":8123")
  -backend string
        后端服务地址
  -debug
        是否打印debug信息
  -tcp.keepalive int
        保持连接的时间间隔(单位秒) (default 15)
  -tcp.timeout int
        连接超时时间(单位秒) (default 5)
  -tls.cert string
        tls证书文件路径(pem格式)
  -tls.key string
        tls密钥文件路径(pem格式)

```

```sh
./tcpproxy \
  -addr=0.0.0.0:8123 \
  -backend=127.0.0.1:8124 \
  -tls.cert=../generate_keys/cert.pem \
  -tls.key=../generate_keys/key.pem \
  -debug
```
