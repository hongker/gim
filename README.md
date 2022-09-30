## GIM
golang实现的基于内存的聊天服务

## 架构
![im](./im.png)

## 功能
- 支持私聊与群聊的消息发送
- 支持按会话查询历史消息
- 支持查询群成员

## 特性
- 推送策略：私聊全量推送，群聊按时间间隔推送最新n条数据
- 消息策略：通过SortedSet数据结构实现消息存储，按发送时间排序

## 启动服务
```
go run cmd/main.go
```

- 查看帮助   
```
> go run cmd/main.go --help
NAME:
   gim - simple and fast im service

USAGE:
   gim [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-address value, --http value          Set http server bind address (default: ":8081")
   --gateway-address value, --ws value        Set websocket server bind address (default: ":8080")
   --gateway-worker value, --ws-worker value  Set websocket server worker number (default: 8)
   --help, -h                                 show help (default: false)
   --profiling-enabled, --profiling           Set pprof switch (default: false)
   --trace-header value, --trace value        Set trace header (default: "trace")
   --version, -v                              print the version (default: false)

```

## 选项设计
```
--enable-profiling: 是否开启pprof
--gateway-address: websocket网关地址
--api-address: http接口地址
--message-protocol: 消息协议(json,proto)
--enable-message-sequence: 是否开启消息sequence
--message-push-count: 每次推送消息条数
--message-push-duration: 每次推送时间间隔
--message-push-concurrency-limit: 消息推送并发上限
--storage: 存储类型(memory, redis)
```
## 模块设计
```

```

## 连接测试
```
# install websocket client
npm install -g wscat

# connect
wscat.cmd -c ws://127.0.0.1:8080 # windows
wscat -c ws://127.0.0.1:8080 # linux

# login
{"op":101,"body":"{\"id\":\"1001\"}"}

# send heartbeat
{"op":105}

# send private message
{"op":201,"body":"{\"content\":\"test\",\"type\":2,\"category\":\"text\",\"target_id\":\"1002\"}"}

# join chatroom
{"op":301,"body":"{\"id\":\"9999\"}"}

# send chatroom message
{"op":201,"body":"{\"content\":\"test\",\"type\":1,\"category\":\"text\",\"target_id\":\"9999\"}"}


# list session
{"op":205}

# query history message
{"op":203,"body":"{\"session_id\":\"1:9999\"}"}
```