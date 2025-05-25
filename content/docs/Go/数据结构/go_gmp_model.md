# Go 运行时调度：G‑P‑M 模型全解

> **适用版本**：Go 1.22（设计自 1.2 奠基，1.14 搭配异步抢占形成当前形态）  
> **核心概念**：Goroutine (`G`)、Processor (`P`)、Machine (`M`)  
> **关键词**：Work‑stealing、Run Queue、Syscall Block、Netpoll、Async Preempt

---

## 1 · 为什么要有 G‑P‑M？

| 目标 | 问题 | 解决方案 |
| ---- | ---- | -------- |
| **高并发** | Goroutine 数可达百万；稀缺资源是核数 | 引入 **P** 承载调度上下文 |
| **低成本切换** | OS 线程切换昂贵（>1 µs） | 在用户态切换 G (`~30 ns`) |
| **充分利用多核** | 全局锁 & 线程竞争严重 | **Work‑stealing** + 局部队列 |
| **隐藏阻塞** | 系统调用/网络 IO 可能卡核 | **M <=> P 解绑**，阻塞不拖慢调度 |

---

## 2 · 三大角色

### 2.1 `G` — Goroutine 控制块

```go
type g struct {
    stack   stack    // [lo, hi)
    sched   gobuf    // 保存寄存器快照
    status  uint32   // _Grunnable, _Grunning, ...
    m       *m       // 正在执行该 g 的 M
    stackguard0, stackguard1 uintptr // 抢占 & 缺栈检查
    // ...省略
}
```

- **生命周期**：`_Gidle → _Grunnable → _Grunning → _Gwaiting/_Gdead …`  
- **特殊 g**：`g0` (系统栈)、`gsignal` (信号栈)

### 2.2 `P` — Processor (逻辑 CPU)

| 字段 | 含义 |
| ---- | ---- |
| `runq [256]g` | 每 P 本地环形队列 |
| `runqhead/runqtail` | 无锁环队列指针 |
| `schedtick` | 调度计数；抢占门控 |
| `mcache` | 本地对象缓存 (alloc fast‑path) |
| `gcAssistTime` | GC 助攻计时 |

`P` 的数量 = **`GOMAXPROCS`**，决定并发度上限。

### 2.3 `M` — Machine (绑定 OS 线程)

| 字段 | 说明 |
| ---- | ---- |
| `g0` | 系统栈；执行调度、GC |
| `curg` | 当前运行的 G |
| `p` | 绑定的 `P` (可空) |
| `park` | 当无 P 可用时睡眠的 semaphore |

- `M` 可 **多于** `P`（阻塞系统调用需额外线程）。

---

## 3 · 调度循环核心流程

```text
┌─ G 进入 runnable ─┐
│  newproc/go func  │
└───────┬───────────┘
        ▼
 P.runq ← push
        ▼                no local G ?
┌── execute G ──┐  ──► steal from other P ─┐
│  switch to g  │                         │
└───► run g     │                         ▼
        │                global runq ?
        │ syscall/block  ───────────────► steal from global
        ▼
  hand P back to sched  
```

### 3.1 本地优先 & Work‑stealing

1. `M` 从 `P.runq` Pop G，若空：  
   1. 从全局队列 `sched.runq` 抢 0.5 批量  
   2. 随机挑一半 `P` 尝试 steal 1/2 G  
2. 新 `G` 创建时优先入 **当前 P.runq**；若满再入全局。

### 3.2 系统调用处理

- **进入**：`g` 进入 `_Gsyscall`；`M` 放回 `P`，自身进入 syscall 队列。  
- **返回**：唤醒休眠 `M` 或新建 `M` 以接管空闲 `P`。

### 3.3 网络 Poller & Timer

- `netpoll` 轮询 epoll/kqueue；就绪 `G` 直接放入当前 `P` 或全局队列。  
- `timer` 到期后同样入队。

### 3.4 抢占

| 方式 | 版本 | 触发 |
| ---- | ---- | ---- |
| **协作抢占** | 全版本 | 函数序言 `stackguard0` 检查 |
| **异步抢占** | ≥1.14 | 向线程发送信号，中断到 `g0` |

---

## 4 · 状态机一览

```text
           +---------+
           | _Gidle  |
           +----+----+
                |
          new G | acquire P
                v
+---------+  runnable   +-----------+
| _Gdead  |-----------> | _Grunnable|
+---------+             +-----+-----+
                              |
                         schedule
                              v
                        +-----+-----+
                        | _Grunning |
                        +--+---+----+
                           |   |
               syscall/block|   |time slice
                           |   |
                    +------+------+----+
                    | _Gsyscall  | _Gwaiting
                    +------------+-------+
```

---

## 5 · 关键源码导航

| 文件 | 关注函数 / 结构 | 说明 |
| ---- | -------------- | ---- |
| `runtime/proc.go` | `schedule`, `execute`, `newproc`, `goexit` | 调度主循环与 G 创建、销毁 |
| `runtime/runtime2.go` | `type g`, `type m`, `type p` | 三大核心数据结构 |
| `runtime/asm_*.s` | `goschedguarded`, `asminit` | 汇编入口、寄存器保存 |
| `runtime/netpoll_kqueue.go` / `_epoll.go` | `netpoll` | I/O 多路复用集成 |

---

## 6 · 调优与观测

| 场景 | 指标 / 工具 | 建议 |
| ---- | ----------- | ---- |
| **调度瓶颈** | `go tool trace` → “Scheduler Latency” | 观察 `Proc Idle` vs `Runnable` |
| **过量线程** | `M->P` Ratio (runtime metrics) | 频繁 syscall → 优化 I/O 模型 |
| **栈拷贝高** | pprof `gcsamples-bytes` | 避免深递归 & 大局部变量 |
| **抢占延迟** | trace `GoSched/GoPreempt` | `runtime/debug.SetMaxThreads(…)` |

---

## 7 · 常见问题排查

| 症状 | 可能原因 | 排查 |
| ---- | -------- | ---- |
| **CPU 利用率低** | GOMAXPROCS < CPU 或 goroutine 阻塞多 | `runtime.GOMAXPROCS(0)` / trace |
| **高延迟尖刺** | 系统调用阻塞；抢占不及时 | `sysmon` 日志 / `schedtrace=1000` |
| **线程爆炸** | cgo 或反复陷入阻塞 | `ps -L` / `debug.SetMaxThreads` |

---

## 8 · 版本演进亮点

| 版本 | 改进 | 影响 |
| ---- | ---- | ---- |
| 1.2 | 引入 P、本地 runq | 大幅提升并发度 |
| 1.10 | netpoll 升级 & timer 分层 | 定时器效率↑ |
| 1.14 | **Async Preempt** | 抢占延迟从 ms → µs |
| 1.20 | 调度器 trace 重写 | 可视化更友好 |
| 1.22 | Goroutine Aware Alloc (mcache) | 分配快‑path 进一步本地化 |

---

## 9 · 实战 Tips

1. **优先本地化**：避免在热点循环中创建/销毁 goroutine，可用对象池或 channel 缓冲复用。  
2. **合理 GOMAXPROCS**：I/O 型服务可略小于 CPU 核；计算型贴合物理核数。  
3. **善用 `go tool trace`**：调度可视化最直接。  
4. **避免长阻塞**：数据库/网络调用加入超时；必要时用 semaphore 控并发。  
5. **使用 `GODEBUG=schedtrace=1000,scheddetail=1`** 观测实时 Runq 长度。

---

### 一句话总结

> **M 执行 G，P 提供执行权，runq + steal 让 goroutine 在多核间均衡，高并发低开销。**

---

如需对 **trace 火焰图、schedtrace 日志** 做深入解析，或定位调度延迟尖刺，欢迎将数据发给我，一起逐帧分析。
