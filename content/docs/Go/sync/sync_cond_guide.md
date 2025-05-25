# Go `sync.Cond` 深度指南

> **适用版本**：Go 1.22  
> **核心定位**：在并发环境中 **等待共享状态变为“满足条件”** 并进行一对一或一对多唤醒。  
> 适用于生产者‑消费者、有界资源池、栅栏同步、事件广播等场景。

---

## 1 · 基本概念

| 组件 | 作用 |
|------|------|
| **Locker (`L`)** | `*sync.Mutex` / `*sync.RWMutex`；`Cond` 自身不带锁 |
| **Wait 队列** | 调用 `Wait()` 的 goroutine 在此休眠（内部用 `runtime_Semacquire`) |

### 核心 API

```go
type Cond struct {
    L Locker          // 必须先 Lock
    // 私有字段
}

func NewCond(l Locker) *Cond
func (c *Cond) Wait()      // 释放锁并阻塞，唤醒后自动重新上锁
func (c *Cond) Signal()    // 唤醒 1 个等待者
func (c *Cond) Broadcast() // 唤醒所有等待者
```

> **规则**：`Signal/Broadcast/Wait` **都必须在持锁状态下调用**，否则运行时 panic。

---

## 2 · 用法模板（Loop‑Wait 习惯用法）

```go
c.L.Lock()
for !condition() {  // **** 必须 for ****
    c.Wait()        // Wait 内部 Unlock → 等待 → 重新 Lock
}
// 条件满足，安全访问/修改共享数据
c.L.Unlock()
```

* 为什么 **for** 而非 **if**？  
  * 可能被 *伪唤醒* 或被 `Broadcast` 唤醒但条件已被其他 goroutine 改变。  
  * 可抵御 “惊群” 竞争：只有条件满足者才真正进入临界区。

---

## 3 · 典型场景与示例

### 3.1 有界队列（生产者‑消费者）

```go
type BoundedQ struct {
    mu              sync.Mutex
    notEmpty, notFull *sync.Cond
    buf             []int
    head, tail, size int
}

func NewBoundedQ(cap int) *BoundedQ {
    q := &BoundedQ{buf: make([]int, cap)}
    q.notEmpty = sync.NewCond(&q.mu)
    q.notFull  = sync.NewCond(&q.mu)
    return q
}

func (q *BoundedQ) Enq(v int) {
    q.mu.Lock()
    for q.size == len(q.buf) { q.notFull.Wait() }
    q.buf[q.tail] = v
    q.tail, q.size = (q.tail+1)%len(q.buf), q.size+1
    q.notEmpty.Signal()  // 唤醒一个消费者
    q.mu.Unlock()
}

func (q *BoundedQ) Deq() int {
    q.mu.Lock()
    for q.size == 0 { q.notEmpty.Wait() }
    v := q.buf[q.head]
    q.head, q.size = (q.head+1)%len(q.buf), q.size-1
    q.notFull.Signal()   // 唤醒一个生产者
    q.mu.Unlock()
    return v
}
```

### 3.2 栅栏（Barrier）同步

```go
type Barrier struct {
    mu      sync.Mutex
    arrived int
    total   int
    cond    *sync.Cond
}

func NewBarrier(total int) *Barrier {
    b := &Barrier{total: total}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Wait() {
    b.mu.Lock()
    b.arrived++
    if b.arrived == b.total {
        b.arrived = 0         // 重用
        b.cond.Broadcast()    // 唤醒全部
    } else {
        for b.arrived != 0 {  // 等待其他 goroutine
            b.cond.Wait()
        }
    }
    b.mu.Unlock()
}
```

### 3.3 事件广播（配置热加载）

```go
var (
    mu      sync.Mutex
    updated bool
    cond    = sync.NewCond(&mu)
)

func Watcher() {
    mu.Lock()
    for !updated {
        cond.Wait()
    }
    // 读取更新后的全局配置
    mu.Unlock()
}

func ReloadConfig() {
    mu.Lock()
    // ...更新配置...
    updated = true
    cond.Broadcast()  // 通知所有 Watcher
    mu.Unlock()
}
```

---

## 4 · Cond vs Channel

| 对比点 | `sync.Cond` | `chan` |
|--------|-------------|--------|
| **语义** | 条件满足后唤醒等待者 | 数据/信号传递 |
| **多等待者广播** | `Broadcast()` 直接支持 | 需 `close(chan)` 或 fan‑out |
| **复合条件** | 任意逻辑 `for !f()` | 复杂 select/多 chan |
| **锁集成** | 必须加锁，适合状态+互斥 | 无锁，自然阻塞 |
| **易用性** | 较复杂，易死锁 | 语法简单，select 强大 |

---

## 5 · 常见坑 & 调试

| 坑 | 说明 |
|----|------|
| **忘记 for 循环** | 伪唤醒或被抢状态导致错误 | 
| **Signal/Broadcast 未持锁** | 运行时会 panic |
| **Unlock 两次** | `Wait()` 返回时已持锁，别再 Unlock 多次 |
| **饥饿** | 大批 Wait 者被广播唤醒抢锁；可分批 Signal |
| **难以超时/取消** | 需外部 `time.After` + goroutine 手动 `Signal` |

**调试技巧**  
* `go test -race` 发现锁竞争  
* `GODEBUG=condout=1`（Go 1.20+）打印 Wait/Signal 关系  
* `go tool trace` 查看阻塞时间

---

## 6 · 性能考虑

* **无等待快路径**：仅加解锁，成本低于 channel select。  
* **大量 Wait 者惊群**：`Broadcast` 后抢锁，可能抖动；必要时 `Signal` + 条件变量。  
* **热点场景**：Lock 区域务必保持极小临界区，避免长阻塞。

---

### 一句话总结

> **`sync.Cond` 是 Go 的“条件变量”：等待共享状态变化，并通过 `Signal/Broadcast` 精准唤醒，占位小、性能高，但必须搭配锁和循环检查条件。**

---

*(完)*