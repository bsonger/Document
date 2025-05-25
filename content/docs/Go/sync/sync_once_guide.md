# Go `sync.Once` 完全指南

> **适用版本**：Go 1.22  
> **核心定位**：并发环境下 **确保某个初始化代码只执行一次**，后续调用立即返回。

---

## 1 · API 速览

```go
type Once struct { /* 非导出 */ }

func (o *Once) Do(f func())   // 调用 f，一生仅一次
func (o *Once) Reset()        // Go 1.21+ 让 Once 可重用
```

| 方法 | 语义 | 关键点 |
|------|------|--------|
| `Do` | *只调用一次* | 其他并发 goroutine 会阻塞直到第一次执行完 |
| `Reset` | 重置执行状态 | 需要外部确保没有并发 `Do`，否则竞态 |

---

## 2 · 内部实现简述

```go
type Once struct {
    done uint32   // 0 未执行，1 已执行
    m    Mutex
}
```

1. **fast‑path**：`atomic.LoadUint32(&o.done)`，为 1 时直接返回。  
2. **slow‑path**：拿锁，二次检查 `done`，再调用 `f()`；成功后 `atomic.StoreUint32(&o.done, 1)`。  
3. **panic 行为**：若 `f()` **panic**，`done` 仍被置 1，后续不会再执行。  

---

## 3 · 典型使用场景

| 场景 | 说明 |
|------|------|
| **懒加载单例** | DB 连接、配置文件、全局缓存 |
| **只注册一次的回调** | Prometheus 指标、HTTP 路由 |
| **高并发多读初始化** | 函数级 static 变量、昂贵运算缓存 |
| **Idempotent 关闭** | `sync.Once` + `close()` 防止重复关闭通道或资源 |

---

## 4 · 示例代码

### 4.1 数据库懒初始化

```go
package db

import (
    "database/sql"
    "sync"

    _ "github.com/go-sql-driver/mysql"
)

var (
    once sync.Once
    conn *sql.DB
)

func Get() *sql.DB {
    once.Do(func() {
        d, err := sql.Open("mysql", "dsn...")
        if err != nil { panic(err) }
        d.SetMaxOpenConns(20)
        conn = d
    })
    return conn
}
```

### 4.2 安全关闭通道

```go
var (
    closeOnce sync.Once
    ch        = make(chan struct{})
)

func SafeClose() {
    closeOnce.Do(func() { close(ch) })
}
```

### 4.3 可重置的 Once（Go 1.21+）

```go
var o sync.Once

func Init() {
    o.Do(func() { fmt.Println("init") })
}

func Reinit() {
    o.Reset()  // 必须确保外部无并发
}
```

---

## 5 · 常见陷阱

| 陷阱 | 说明 |
|------|------|
| **panic 不重试** | `f()` panic 后 `done=1`，后续 `Do` 会直接返回，可能导致系统缺少初始化。 |
| **闭包引用变量** | 若闭包里访问外部可变状态，需保证可见性和一致性。 |
| **误用 Reset** | 只能在确定无并发 `Do` 的情况下调用，否则数据竞争。 |

---

## 6 · 性能 & 开销

* **fast‑path** 只有一次原子读：开销 ~2–3 ns  
* **slow‑path** 仅第一次执行，用锁保护。  
* 相比双重检查锁（DCL），`sync.Once` 代码可读且正确性经验证。

---

## 7 · 版本演进

| 版本 | 变更 | 影响 |
|------|------|------|
| 1.0  | `sync.Once` 引入 | 保证单次调用 |
| 1.12 | fast‑path 原子优化 | 在多核下吞吐明显提升 |
| 1.21 | `(*Once).Reset()` | 支持重用 Once |

---

### 一句话总结

> **`sync.Once` = 线程安全且零开销 fast‑path 的“只执行一次”工具；善用它实现懒初始化和幂等关闭。**

---

*(完)*