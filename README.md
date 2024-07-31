# MemoryDB

MemoryDB 是一个简单的内存数据库，用于存储键值对并支持过期管理。它提供了基本的 `SET`、`GET`、`DEL` 和 `EXPIRE` 操作，并具有定期清理过期键的功能
适合单服务且为了节约成本未使用 Redis 的小公司

## 功能

- 设置键值对，并支持设置过期时间
- 获取键值对，处理键不存在或已过期的情况
- 删除键值对
- 设置键值对的过期时间
- 定期清理过期键
- 关闭数据库，停止清理任务并清空存储数据

## 安装

使用 `go get` 下载并安装 MemoryDB：

```sh
go get github.com/nuominmin/memorydb
```
