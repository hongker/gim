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

## 启动
```
cd cmd
go build server.go
```

- 启动服务   
```
>./gim run --help
NAME:
   gim run - run service

USAGE:
   gim run [command options] [arguments...]

OPTIONS:
   --config FILE, -c FILE     Load configuration from FILE (default: "./app.yaml")
   --debug                    Set debug mode (default: false)
   --limit value, -l value    Set max number of session history messages (default: 10000)
   --port value, -p value     Set tcp port (default: 8080)
   --storage value, -s value  Set storage (default: "memory")

```

- 查看版本号
```
./gim --version
```
