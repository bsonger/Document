---
title: "introduce"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---

# 🚀 Cilium 简介、架构与作用

## 📌 一、Cilium 是什么？

Cilium 是一个基于 eBPF 技术构建的云原生网络和安全平台，它作为 Kubernetes 的 CNI 插件，提供了高性能的容器网络连接、网络安全策略、流量可观测性以及 Service Mesh 能力。

### ✨ 核心特点：
- 🧠 基于 Linux eBPF 技术，性能优异、无侵入
- 🔐 支持 L3/L4/L7 网络策略（支持 HTTP/gRPC/Kafka 等）
- ♻️ 无需 Sidecar 的 Service Mesh 能力
- 👀 内建 Hubble 实现网络流量可视化与追踪
- 🌐 支持多集群互联（ClusterMesh）

---

## 🏗️ 二、Cilium 架构

Cilium 架构主要由以下几个组件组成：

### 1️⃣ Cilium Agent
- 每个 Node 上运行
- 管理和编译 eBPF 程序，注入到内核
- 处理服务发现、负载均衡、策略同步等

### 2️⃣ eBPF 程序（运行在内核）
- 挂载在 XDP、TC、Socket 等钩子点上
- 实现数据包的解析、过滤、重定向、转发
- 性能高，运行在内核态

### 3️⃣ Hubble
- 基于 eBPF 的可观测性平台
- 收集网络流量事件、DNS、HTTP 等信息
- 提供 CLI 与 UI 工具，支持 Grafana 集成

### 4️⃣ Cilium CLI
- 管理工具，支持安装、诊断、策略测试等

### 5️⃣ Cilium Operator
- 控制平面组件，处理 CRD、Kubernetes 资源变更

### 6️⃣ ClusterMesh（多集群）
- 实现多个 Kubernetes 集群之间的 Pod 互联
- 同步服务、策略、身份信息

---

## 🛠️ 三、Cilium 的主要作用

| 📌 场景 | 描述 |
|--------|------|
| **🔗 CNI 插件** | 提供容器网络连接，支持 IPv4/IPv6、Overlay、BGP 等 |
| **🛡️ 网络安全策略** | 支持从 L3（IP）到 L7（HTTP Path）的多层访问控制 |
| **🧩 Service Mesh** | 支持无 Sidecar 模式，支持 mTLS、L7 策略、流量控制 |
| **🔍 可观测性** | Hubble 提供流量追踪、指标、路径分析等功能 |
| **🚪 Egress Gateway** | 精细控制 Pod 的出口访问地址与 IP |
| **🌐 多集群通信** | 支持 ClusterMesh，多个集群共享网络与策略 |

---

## 📈 四、适用场景

- 🚀 构建高性能 Kubernetes 网络平台
- 🧵 对服务间通信进行精细化控制（如 L7 HTTP 策略）
- 🧼 替代传统 Service Mesh 实现零 Sidecar 化
- 🛡️ 构建零信任网络架构（mTLS + 策略）
- 🔄 跨集群网络管理
- 🧪 网络观测、调试与安全审计

---

## 🧾 五、总结

Cilium 是一个集网络连接、安全策略、可视化与服务网格功能于一体的现代化平台，特别适合对性能、安全、可观测性有较高要求的生产集群环境。凭借 eBPF 技术优势，Cilium 正逐渐成为下一代 Kubernetes 网络标准的有力竞争者。
