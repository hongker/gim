## GIM
基于内存的IM聊天系统

## 架构设计
- gate: 长链接网关，负责维护客户端的连接、数据收发。不参杂业务逻辑，避免因业务迭代而重启。
- logic: 业务逻辑层，负责处理由网关转发的客户端请求，以及主动推送数据到客户端。

## 协议设计
报文长度(x) = 4 + 2 + 2 + n
```
-----------------------------------
| 报文长度 | 版本 | 操作类型 | 数据项 |
-----------------------------------
|   4位   | 2位  |  2位    |  n位  |
-----------------------------------
```

## 功能设计
- 支持私聊
- 支持群聊

## 性能测试