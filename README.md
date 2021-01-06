# 一个简单的 tcp 代理

```sh
Usage of tcpproxy:
  -addr string
        代理的监听地址 (default ":8123")
  -backend string
        后端服务地址
  -cert string
        tls证书文件路径(pem)
  -debug
        是否打印debug信息
  -key string
        tls密钥文件路径(pem)

```

```sh
./tcpproxy -backend=127.0.0.1:8124 -cert=../generate_keys/cert.pem -key=../generate_keys/key.pem -debug
```
