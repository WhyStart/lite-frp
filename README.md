# lite-frp
lite-frp是一个Go语言编写的一个内网穿透工具

lite-frp实现了基于TCP的私有协议
# 后续计划
* 实现数据的传输加密
* 支持多种协议
# 配置文件说明
conf/client.ini
```ini
[common]
server_addr = 127.0.0.1 // 服务端ip
server_port = 6000 // 服务端与客户端通信端口
[local]
type = tcp // 传输协议
sk = abcdefg // 密钥
local_ip = 127.0.0.1 // 客户端ip
local_port = 8090 // 客户端端口
remote_port = 9090 // 客户端外网端口
```
conf/client.ini
```ini
[common]
bind_port = 6000 // 服务端端口
type = tcp // 传输协议
key = abcdefg // 密钥
```

