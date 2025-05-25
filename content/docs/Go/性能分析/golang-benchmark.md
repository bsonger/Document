---
title: "benchmark"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---
# Go 基准测试（Benchmark）原理与实践

> 📈 **go test -bench** 能让你在纳秒级别洞察代码性能。本指南汇总了从运行原理到实战技巧的完整知识点。

## 1 — 核心概念速览

| 关键对象 | 作用 |
|---------|------|
| `go test -bench` | 启动基准测试运行器（benchmark harness）。 |
| `func BenchmarkXxx(b *testing.B)` | 声明一个基准测试函数，名称需以 `Benchmark` 前缀开头。 |
| `b.N` | 运行器自动决定的循环次数；被测代码需放入 `for i := 0; i < b.N; i++ { … }`。 |
| 计时/统计 | 默认统计 **wall‑clock** 时间（`ns/op`）；使用 `-benchmem` 可额外输出 `B/op` 与 `allocs/op`。 |
| 常用辅助 API | `b.ResetTimer()`、`b.StopTimer()/b.StartTimer()`、`b.ReportAllocs()`、`b.SetBytes(n)`、`b.RunParallel()` |

---

## 2 — Benchmark 运行原理

1. **编译阶段**  
   - `go test` 会为测试与基准生成临时 `main` 包并编译。

2. **探测合适的迭代次数 (`b.N`)**  
   - 初始化 `N = 1`；如果一次耗时 \< 0.5 s，就按指数级放大（×20、×100 …）。  
   - 目标：单次基准 **≥ 1 s**（可通过 `-benchtime` 调整），兼顾统计显著性与总时长。

3. **精准计时**  
   - 默认计时整个基准函数；调用 `StopTimer` 后暂停，`StartTimer` 恢复。  
   - 任何 GC、系统调用也被计入；避免在计时区内执行 I/O、网络等操作。

4. **结果输出**  
   示例：  
   ```text
   BenchmarkNaive-8      1234567      100.0 ns/op
   ```  
   `-benchmem` 追加每次迭代的分配字节数与次数。若调用 `SetBytes(n)`，还会显示 **MB/s**。

5. **并行模式**  
   `b.RunParallel(func(pb *testing.PB){ … })` 会根据 `GOMAXPROCS` 启动多 goroutine，并通过 `pb.Next()` 控制循环。

---

## 3 — 如何书写 & 运行基准

```go
// concat_test.go
package concat

import (
    "strings"
    "testing"
)

func naive(s1, s2 string) string { return s1 + s2 }

func buffered(s1, s2 string) string {
    var buf strings.Builder
    buf.WriteString(s1)
    buf.WriteString(s2)
    return buf.String()
}

func BenchmarkNaive(b *testing.B) {
    s1, s2 := "foo", "bar"
    for i := 0; i < b.N; i++ {
        _ = naive(s1, s2)
    }
}

func BenchmarkBuffered(b *testing.B) {
    s1, s2 := "foo", "bar"
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = buffered(s1, s2)
    }
}
```

### 3.1 运行命令

```bash
# 仅跑基准（跳过单元测试）
go test -run=^$ -bench=. ./...

# 查看内存分配
go test -run=^$ -bench=. -benchmem ./...

# 自定义持续时间（至少 3 s）
go test -run=^$ -bench=. -benchtime=3s ./...

# 生成 CPU / 内存剖面
go test -run=^$ -bench=. -cpuprofile cpu.out -memprofile mem.out
go tool pprof cpu.out
```

---

## 4 — 常用技巧与坑

| 做法 | 理由 |
|------|------|
| 将准备工作放在循环外或使用 `StopTimer`/`StartTimer` | 避免把非核心逻辑计入基准时间 |
| 使用结果变量（如 `var sink any`） | 防止编译器死代码消除 (DCE) |
| 固定输入规模 | 提高可重复性；多规模可用 `b.Run` |
| 注意共享状态 & GC | 共享 map、随机数种子等会放大噪声 |
| 在实体机上运行 | VM、笔记本省电模式可能抖动 |
| 结合 `go test -count=N` | 多次运行取平均/中位提高置信度 |

---

## 5 — 进阶能力

### 5.1 子基准

```go
func BenchmarkSizes(b *testing.B) {
    for _, size := range []int{1 << 10, 1 << 20, 1 << 24} {
        b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
            data := make([]byte, size)
            for i := 0; i < b.N; i++ {
                _ = hash(data)
            }
        })
    }
}
```

### 5.2 并行吞吐测试

```go
func BenchmarkPool(b *testing.B) {
    pool := sync.Pool{New: func() any { return make([]byte, 1024) }}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            buf := pool.Get().([]byte)
            // ... 使用 buf ...
            pool.Put(buf)
        }
    })
}
```

### 5.3 结果可视化

使用官方工具 `benchstat` / `benchcmp` 对比两次基准差异：

```bash
go test -bench=. > old.txt
# 修改代码…
go test -bench=. > new.txt
benchstat old.txt new.txt
```

---

## 6 — 小结

- 基准测试集成于 **`testing` 包**，无需额外依赖。  
- 把 **真正欲测的语句** 放进 `for i < b.N`；其余逻辑用 `StopTimer` 隔离。  
- `-benchmem`、`pprof` 与 `benchstat` 可帮助你定位瓶颈与分配热点。  
- **基准 ≠ 压测**：它关注微观性能；系统级负载需另行压测。  

---
