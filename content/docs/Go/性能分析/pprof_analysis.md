# Go `pprof` 性能分析实践手册

> **适用版本**：Go 1.20 – Go 1.23（2025 年最新版）  
> **阅读对象**：希望系统掌握 `pprof` 并落地到生产的 Go 开发 / SRE

---

## 目录
1. 概览：`pprof` 能解决什么问题  
2. Profile 类型与采样原理  
3. 离线采集：基准测试 & CLI  
4. 在线采集：`net/http/pprof` & 远程抓取  
5. 可视化：命令行、FlameGraph、Speedscope  
6. 案例一：CPU Hotspot 定位与优化  
7. 案例二：内存泄漏排查  
8. 案例三：阻塞 / Mutex 竞争分析  
9. 持续 Profiling：Parca / Pyroscope / Grafana Pyroscope Cloud  
10. PGO：Profile-Guided Optimization 流程  
11. 自动化工作流与 CI 集成  
12. 常见误区 & 调参清单  
13. 参考链接

---

## 1. 概览：`pprof` 能解决什么问题

- **定位 CPU 热点**：采样型 Profile，找到最耗时的函数或调用链  
- **排查内存泄漏**：对比多份堆快照，确定分配点  
- **锁竞争 / 阻塞**：查看被 `sync.Mutex` 或系统调用阻塞的时间  
- **Goroutine 爆炸**：快照当前所有 goroutine 栈，分析创建源

`pprof` 不会直接告诉你“怎么改”，但能精准指向**瓶颈位置**。随后通过代码 / 参数 / 架构层面优化并 **benchstat** 验证。

---

## 2. Profile 类型与采样原理

| Profile              | 默认采样      | 主要字段         |
|----------------------|--------------|------------------|
| **CPU**              | 100 Hz 时钟中断 | `flat`, `cum`    |
| **Heap (allocs/inuse)`| 每分配 512 KiB | `alloc_objects` / `inuse_space` |
| **Block**            | 100 Hz       | `contentions_ns` |
| **Mutex**            | 50 Hz        | `contentions`, `delay` |
| **Goroutine**        | 快照         | `stack`          |

- Go 运行时在热点处插入采样回调，成本≈1–2 % CPU。  
- `Heap` 采样率可通过环境变量 `GODEBUG=memprofilerate=N` 调节。  

---

## 3. 离线采集：基准测试 & CLI

```bash
# 编译基准并输出 cpu.prof
go test -run=^$ -bench=. -benchtime=10s -cpuprofile cpu.prof ./pkg

# 生成内存 profile
GODEBUG=memprofilerate=131072 go test -run=^$ -bench=. -memprofile mem.prof ./pkg
```

> *Tips*  
> - `-run=^$` 跳过单元测试，仅跑基准  
> - `-benchmem` 报告 alloc/op，配合 `mem.prof` 效果更佳

---

## 4. 在线采集：`net/http/pprof`

```go
import _ "net/http/pprof"

func main() {
    go http.ListenAndServe(":6060", nil)
    ...
}
```

### 常用抓取命令

| URL | 说明 |
|-----|-----|
| `/debug/pprof/profile?seconds=30` | CPU 30 s |
| `/debug/pprof/heap` | 堆快照 |
| `/debug/pprof/mutex?debug=1` | Mutex |
| `/debug/pprof/block?debug=1` | Block |

确保只对内网或 mTLS 暴露； kubernetes 可用 `kubectl port-forward` 临时拉通。

---

## 5. 可视化

### 5.1 CLI

```bash
go tool pprof -top cpu.prof
go tool pprof -list 'MyHotFunc' cpu.prof
```

### 5.2 内置 Web

```bash
go tool pprof -http=:8080 cpu.prof
```

### 5.3 Speedscope

```bash
go tool pprof -proto cpu.prof > cpu.pb.gz
# 将 protobuf 文件拖到 https://www.speedscope.app
```

Speedscope 支持时间轴视图，适合分析并行 / I/O 阻塞。

---

## 6. 案例一：CPU Hotspot

1. **症状**：QPS 达到 5k 时 CPU 占用 900 %  
2. **采集**：30 秒 CPU profile  
3. **分析**：`top` 发现 `regexp.(*Regexp).FindAllIndex` flat=45 %  
4. **优化**：  
   - 复用编译后的正则  
   - 替换为 `strings.Contains` + 预解析  
5. **回归**：benchstat 显示 **latency ↓30 % / CPU ↓25 %**

---

## 7. 案例二：内存泄漏

```bash
# 周期抓取堆
go tool pprof http://app:6060/debug/pprof/heap -o heap1.prof
sleep 120
go tool pprof http://app:6060/debug/pprof/heap -o heap2.prof
# Diff
go tool pprof -base heap1.prof heap2.prof
```

若差异集中在 `bytes.makeSlice`，通常是切片无限扩容；检查 `cap` 预估或扩容策略。

---

## 8. 案例三：阻塞 / Mutex

- **Block profile**：定位网络 / DB 延迟  
- **Mutex profile**：出现 `runtime.mapassign` 热点 → map 写争用；可换 `sync.Map` 或 sharding

Go 1.21+ 支持 `runtime/pprof.Labels("tenant", id)`，在火焰图中按标签聚合，快速识别哪条租户耗时。

---

## 9. 持续 Profiling

| 方案 | 特点 |
|------|------|
| **Parca / Parca Agent** | 云原生，eBPF + DWARF，无需 sidecar |
| **Pyroscope** | UI 友好，可与 Grafana Merge |
| **Go 运行时 mprof** | 1.20+ 内置连续 CPU/Heap 采样，低开销 |

在 Kubernetes 集群推荐 DaemonSet 方式安装 Agent，统一上报到集中 UI。

---

## 10. PGO：Profile-Guided Optimization

1. Go 1.20 起支持 `-pgo=`  
2. 流程：  
   ```bash
   go test -c -pgo=./pgo.pprof
   go build -pgo=./pgo.pprof -o app
   ```  
3. PGO 主要优化 **内联 / 分支预测**；对微服务效果 3–10 %。

---

## 11. 自动化工作流

```bash
# Makefile 片段
bench:
    go test -run=^$ -bench=. -benchmem -cpuprofile cpu.prof ./...
    benchstat old.txt new.txt
    go tool pprof -svg cpu.prof > cpu.svg
```

CI 中上传 `cpu.svg` 作为 artifact，配合 GitHub Actions 注释 PR。

---

## 12. 常见误区 & 调参清单

| 误区 | 纠正 |
|------|------|
| 抓一次 profile 就够 | 采样有随机性，至少多跑 3 次 |
| Heap Profile = RSS | 只包含 Go 堆，不含 mmap、Cgo、栈 |
| `pprof` 会拖慢线上 | 采样型开销极低，阻塞 / mutex 建议临时开启 |

**调参**  
- `runtime.SetMutexProfileFraction(10)` – 每 10 次竞争采样一次  
- `runtime.SetBlockProfileRate(1)` – 捕捉所有阻塞（排障后改回 0）

---

## 13. 参考链接

- Go 官方 Blog《Profiling Go Programs》  
- Brendan Gregg FlameGraph 项目  
- Parca.dev / Pyroscope.io / Grafana Cloud Continuous Profiling  
- Go 1.23 Release Notes – PGO & profiler 改进

---

> **建议**：在每个重大版本升级后，重新评估采样率与阈值。Go 运行时及 GC 每年都在进步，历史“玄学调参”可能已过时。
