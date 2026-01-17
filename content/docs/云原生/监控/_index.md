---
title: "监控"
weight: 2
bookCollapseSection: true
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Observability Platform Services

## Alloy
| Service | Description |
|---------|-------------|
| alloy | Alloy 是 Pyroscope 的核心存储与处理组件，接收采集器发送的 profiling 数据，并提供查询和聚合能力。 |

## Grafana
| Service | Description |
|---------|-------------|
| grafana | Grafana 前端可视化界面，用于展示 Prometheus、Loki、Tempo、Pyroscope 等监控数据和指标，支持 Dashboard 配置和告警。 |

## Kubernetes Metrics
| Service | Description |
|---------|-------------|
| kube-state-metrics | 收集 Kubernetes 对象状态（Pod、Deployment、Node 等）指标并暴露给 Prometheus，便于监控和告警。 |
| node-exporter | NodeExporter 在节点上收集主机级别指标，如 CPU、内存、磁盘、网络等，并提供给 Prometheus。 |

## Loki (日志系统)
| Service | Description |
|---------|-------------|
| loki-canary | Loki Canary 测试日志管道是否正常。 |
| loki-chunks-cache | 缓存日志块，加速查询。 |
| loki-compactor | 压缩日志块，减少存储占用。 |
| loki-distributor | 分发日志条目到 ingester，实现负载均衡。 |
| loki-index-gateway | 索引查询服务，将查询请求路由到 ingester。 |
| loki-ingester | 日志存储节点，负责存储日志块。 |
| loki-querier | 查询日志数据，并返回给前端或 API。 |
| loki-query-frontend | 查询前端，路由用户请求到 querier。 |
| loki-query-scheduler | 查询调度器，分发查询请求到 querier。 |
| loki-results-cache | 缓存查询结果，加快重复查询速度。 |

## Prometheus (指标系统)
| Service | Description |
|---------|-------------|
| prometheus-k8s | Prometheus 主实例，采集集群和应用指标。 |
| prometheus-operator | Prometheus Operator 控制器，管理 Prometheus、ServiceMonitor、PodMonitor 等 CRD。 |

## Pyroscope (Profiling 系统)
| Service | Description |
|---------|-------------|
| pyroscope-ad-hoc-profiles | 收集临时 profiling 数据，例如 CPU/内存采样，用于临时调试。 |
| pyroscope-admin | 管理后台，提供 Web UI 和 API，管理用户、Tenant、配置等。 |
| pyroscope-compaction-worker | 压缩和合并原始 profiling 数据，提高存储和查询效率。 |
| pyroscope-distributor | 接收 agent 发送的 profiling 数据，并分发到 segment writer 或 metastore。 |
| pyroscope-metastore | 保存 profiling 索引、应用信息和标签，支撑查询和聚合。 |
| pyroscope-query-backend | 聚合 metastore 和 segment writer 的数据，处理查询请求。 |
| pyroscope-query-frontend | Web 前端展示，调用 query-backend 渲染 flamegraph 和表格。 |
| pyroscope-segment-writer | 存储原始 profiling 数据，按时间段保存 stack trace。 |
| pyroscope-tenant-settings | 多租户配置管理，如采样率、存储策略、权限配置。 |

## Tempo (Tracing 系统)
| Service | Description |
|---------|-------------|
| tempo-compactor | 压缩 Trace 数据块，优化存储并提高查询性能。 |
| tempo-distributor | 分发 Trace 数据到 ingester，实现负载均衡。 |
| tempo-ingester | 存储 Trace 数据块，并提供查询接口。 |
| tempo-memcached | 缓存 Trace 数据，加速查询响应。 |
| tempo-metrics-generator | 生成 Tempo 相关指标数据，用于监控集群状态。 |
| tempo-querier | 查询 Trace 数据，为前端提供接口。 |
| tempo-query-frontend | 前端查询服务，将请求路由到 querier 并返回结果。 |