# Go 垃圾回收（GC）工作原理全景

> **适用版本**：Go 1.22；从 1.5 并发三色标记‐清扫奠基，后续演进始终保持总体框架不变。  
> **关键词**：并发三色、写屏障、Mutator Assist、Scavenger、STW 压缩、Soft Memory Limit。

---

## 1 · 总体架构

| 维度 | 特性 | 说明 |
| ---- | ---- | ---- |
| **算法** | 并发三色 **Mark‑Sweep** | 绝大部分标记、清扫与业务线程并发执行 |
| **停顿** | 两次 **STW** | *Mark Setup* & *Mark Termination*，通常 < 1 ms |
| **并发** | 后台/辅助标记 + 清扫 | 业务线程按“分配比例”主动干 GC，后台 goroutine 持续工作 |
| **内存回收** | Span 清扫 + **Scavenger** | 小对象归还到 per‑P 空闲链，大 Span 退还操作系统 |

---

## 2 · 一次 GC 循环

```text
┌──────────────┐      ┌─────────────────────────────┐
│ STW #1       │      │ STW #2                      │
│ Mark Setup   │──▶──▶│ Mark Termination (黑完灰)    │
└─▲────────────┘      └────────▲────────────────────┘
  │                            │
  │并发写屏障                  │并发 Sweep / Scavenge
  └───► 并发标记 (Mutator Assist + BG marker) ────►┘
```

### 2.1 阶段要点

| 阶段 | 作用 | 典型耗时 (8C/2 GiB) |
| ---- | ---- | ------------------- |
| **Mark Setup (STW)** | 根对象灰化（栈/全局/寄存器）；打开写屏障 | 30–200 µs |
| **并发标记** | 灰 → 黑；Mutator Assist 按分配量帮忙 | 1–100 ms (并发) |
| **Mark Termination (STW)** | 灰队列清空；停写屏障；记录 liveBytes | 50–400 µs |
| **并发清扫 (Sweep)** | 遍历 span，将空闲块归还 mcentral | 1–100 ms (并发) |
| **Scavenger** | 大块空闲页 madvise 回 OS，压 RSS | 持续后台 |

---

## 3 · 三色标记 & 写屏障

### 3.1 颜色含义

| 颜色 | 状态 | 行为 |
| ---- | ---- | ---- |
| **白** | 未访问 | 可能垃圾 |
| **灰** | 已发现未扫描 | 等待扫描 |
| **黑** | 已扫描 | 活对象 |

> **不变式**：**黑对象不能指向白对象**。  
> 写屏障负责在 Mutator 并发写指针时，把 *旧值* 或 *新值* 灰化，维持不变式。

### 3.2 Hybrid Write Barrier

```go
dstOld := *p         // 读取旧指针
* p = newPtr         // 写新指针
shade(dstOld)        // 灰化旧值
```

- 编译器在插桩阶段对 “可能在 GC 时写指针” 的指令包裹 `writebarrierptr`。
- 只有在 *并发标记窗口* 内屏障才做实事，其余时间是空函数 → 99.9% 路径零开销。

---

## 4 · Mutator Assist 机制

- 分配路径会累加 **assist debt**  
  `debt += bytesAllocated * (markWorkRate / allocRate)`  
- 若负债超阈值，则在分配调用栈中 **主动扫描灰对象**，直到偿还债务。  
- 这样实现 **“GC 速度 ≥ 分配速度”** 的自适应节奏，无需频繁全局锁。

---

## 5 · Heap 布局与清扫

| 概念 | 粒度 | 功能 |
| ---- | ---- | ---- |
| **Span** | 8 KiB–4 MiB | 内存管理基本单元，记录对象位图 |
| **Class (size class)** | 8 B–32 KiB 共 67 档 | 小对象按 class 分配，复用缓存 |
| **MCache (每 P)** | 线程/处理器本地空闲链 | 分配无需锁 |
| **MCentral** | 全局空闲链 | Sweep 将 Free span 回收 |

**Sweep**：  
1. 读取 span 位图，黑对象所在 block 标记保留。  
2. 其余 block 归并到 freeList；若整 Span 空，挂到 `mheap.sweepSpansEmpty`。  
3. 分配若无缓存，会触发 On‑Demand Sweep。

---

## 6 · Scavenger (退还物理内存)

- Go 1.17 起 **背景线程 + Bitmap** 追踪 *unused* 页。  
- 达阈值后调用 `madvise(MADV_DONTNEED)` → OS 立刻回收物理页，RSS 下降。  
- 受 `GODEBUG=madvdontneed=1` 控制；默认为开。

---

## 7 · 堆目标 & 节奏控制

```go
targetHeap = liveHeap * (1 + GOGC/100)   // 默认 GOGC=100
```

- 分配达到 `targetHeap` → 触发下一轮 GC。  
- `GOGC`↓ 提前触发、降低峰值内存；`GOGC`↑ 反之但暂停次数减。  
- Go 1.19 引入 **Soft Memory Limit**：若进程 RSS 超阈值，GC 会加倍速率直至降回。

---

## 8 · 观测指标速查

| Metric (Prometheus/runtime) | 解读 |
| --------------------------- | ---- |
| `gc.pause_ns.sum` / `gc.pause_total_ns` | 累计 STW；尾延迟敏感业务建议 < 100 ms/min |
| `gc.cycles.total` / `gc.heap_objects` | GC 频率 & 存活对象趋势 |
| `gc.cpu.ns / process_cpu_seconds_total` | GC CPU 比例；>30% 需优化 |
| `heap_released` vs `heap_idle` | 释放给 OS 的页 vs 未用页；碎片观测 |
| `gc.assist_cpu_fraction` (*1.22*) | Mutator Assist 占用；过高说明分配过猛 |

---

## 9 · 版本演进亮点

| 版本 | 重要变化 | 效果 |
| ---- | -------- | ---- |
| **1.14** | Async preempt；函数中途抢占 | STW 压缩 × 2 |
| **1.17** | 新 Scavenger；Huge Page 感知 | 容器 RSS ‑30% |
| **1.19** | Soft Memory Limit；Pacing 重写 | OOM 风险显著下降 |
| **1.22** | 写屏障路径简化；灰队列优化 | 2 GiB 堆 STW < 1 ms |
| **1.23 (计划)** | 年轻代（Generational）原型 | 进一步削减短命对象成本 |

---

## 10 · 实战建议

1. **先量化**：打开 `runtime/metrics → Prometheus`，确定究竟是暂停、CPU、还是 RSS 瓶颈。  
2. **热点代码审逃逸**：`go build -gcflags=-m` + `benchmem` 量化 **allocs/op**。  
3. **合理调 `GOGC`**：内存紧 → 50–80；CPU 忙 → 150–200；上线前压测找平衡点。  
4. **容器环境**：设置 `GOMEMLIMIT` 或 `GODEBUG=softmemlimit=…`，避免 K8s OOMKill。  
5. **升级 Go 版本**：大版本 GC 改进显著，优先考虑 1.22+。

---

### 一句话记忆

> **“极短 STW + 并发标记 + 自助助攻 + 背景清扫 + 物理页回收” = Go GC 成熟方案。**

---

若你手上有 **gctrace 日志、heap pprof、trace.out** 等文件，或想针对特定场景（延迟、内存、CPU）调优，随时发给我，我们可以逐行解读并给出参数与代码级建议。
