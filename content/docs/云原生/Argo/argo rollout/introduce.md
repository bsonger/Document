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

# 🚀 Argo Rollouts 概述与发布策略

---

## ✅ 什么是 Argo Rollouts？

**Argo Rollouts** 是一个 Kubernetes 控制器，用于实现高级的部署策略，如：

- 蓝绿部署（Blue-Green Deployment）
- 金丝雀发布（Canary Release）
- 实验分析（Analysis Run）
- 支持自动回滚、手动推进、AB 测试等能力

它是对原生 Kubernetes `Deployment` 的增强版，提供更精细的发布控制与观察能力。

---

## 🔄 Deployment vs Rollout

| 特性 | Kubernetes Deployment | Argo Rollouts |
|------|------------------------|----------------|
| 基本滚动更新 | ✅ 有 | ✅ 有 |
| 蓝绿发布 | ❌ 无 | ✅ 支持 |
| 金丝雀发布 | ❌ 无 | ✅ 支持 |
| 手动控制发布阶段 | ❌ 无 | ✅ 支持 |
| 自动分析并决策 | ❌ 无 | ✅ 支持 |
| 版本对比和观察窗口 | ❌ 无 | ✅ 支持 |
| 指标判断（Prometheus 等） | ❌ 无 | ✅ 支持 |
| Web UI 支持 | ❌ 无 | ✅ 有（Argo Rollouts UI 插件） |
| 实验/AB 测试 | ❌ 无 | ✅ 支持 |

---

## 🚦 支持的发布策略

### 1. 金丝雀发布（Canary）

将新版本流量逐步放量，比如 10% -> 30% -> 100%，每个阶段可通过 metrics 分析判断是否推进。

```yaml
strategy:
  canary:
    steps:
    - setWeight: 10
    - pause: {duration: 1m}
    - setWeight: 50
    - pause: {}
```

### 2. 蓝绿部署（Blue-Green）

部署新版本于“预览”环境，手动或自动切换流量。

```yaml
strategy:
  blueGreen:
    activeService: my-app-active
    previewService: my-app-preview
    autoPromotionEnabled: false
```

### 3. 实验分析（Analysis Run）

可插入 Prometheus/Datadog/New Relic 等指标平台，进行自动化判断是否回滚或推进。

```yaml
analysis:
  templates:
    - name: success-rate-check
      metrics:
        - name: error-rate
          provider:
            prometheus:
              query: rate(http_requests_total{status=~"5.."}[1m])
```

---

## 🧩 Argo Rollouts 架构图

![Argo Rollouts Architecture](https://argo-rollouts.readthedocs.io/en/stable/assets/rollouts_architecture.png)

### 关键组件说明：

| 组件 | 说明 |
|------|------|
| Rollout Controller | 控制器核心，替代原生 Deployment 控制器 |
| Rollout CRD | 定义替代 Deployment 的高级对象 |
| AnalysisRun | 执行分析模板、连接外部指标源 |
| Experiment | 并行部署多个版本，用于 AB 测试 |
| UI 插件 | 可集成至 Argo CD 观察 rollout 状态 |

---

## 📚 推荐阅读

- 官方文档：https://argo-rollouts.readthedocs.io/
- GitHub 项目：https://github.com/argoproj/argo-rollouts
- 示例 YAML：https://github.com/argoproj/argo-rollouts/tree/master/examples
