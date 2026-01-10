---
title: "trace"
weight: 3
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: true
---
🔹 Trace 的工作原理

1️⃣ 请求进入系统
•	入口（如 API Gateway）生成一个 Trace ID
•	记录 起始时间

2️⃣ 请求传播
•	每个微服务 创建一个 Span
•	传递 Trace ID（上下文）
•	记录 时间、请求参数、响应时间

3️⃣ 请求完成
•	计算 每个服务的执行时间
•	记录 错误信息
•	上报 Trace 数据到 存储 & 可视化系统

🔹 Trace 数据结构

一个 Trace 由多个 Span 组成，每个 Span 记录：

```json
{
  "trace_id": "abcd-1234",
  "span_id": "efgh-5678",
  "parent_id": "ijkl-9101",
  "service": "user-service",
  "operation": "GET /user/info",
  "start_time": 1710000000000,
  "duration_ms": 30,
  "tags": {
    "http.status_code": 200,
    "db.query": "SELECT * FROM users"
  }
}
```

🚀 总结

✅ Trace 适用于分布式微服务，帮助分析慢请求、调用链

✅ 核心概念：Trace、Span、上下文传播、存储 & 可视化

✅ Jaeger、OpenTelemetry、Tempo 是常见的 Trace 方案

✅ 相比传统日志，Trace 更适合故障排查 & 性能优化
