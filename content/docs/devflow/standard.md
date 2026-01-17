# 云原生平台安全与治理规范（落地化指南）

> 本文档旨在将云原生平台中的安全、治理、SRE、成本、自服务等规范落地，通过自动化手段强制执行，确保平台交付可信、可控、可审计。

---

## 目录

1. [构建与交付](#构建与交付)
2. [安全](#安全)
3. [平台治理](#平台治理)
4. [发布质量 & SRE](#发布质量--sre)
5. [运行时安全](#运行时安全)
6. [成本与资源](#成本与资源)
7. [自服务 & 开发者体验](#自服务--开发者体验)
8. [最小落地集（MVP）](#最小落地集mvp)
9. [附录](#附录)

---

## 构建与交付

| 规范（人话） | 自动化方式 | 强制点 |
|----|----|----|
| 所有镜像必须由平台 CI 构建 | Tekton Pipeline 统一入口 | 禁止手工 push |
| 镜像必须签名 | Tekton + cosign 自动签名 | Admission 校验 |
| 不允许使用长期密钥 | OIDC / SA 临时凭证 | CI 模板 |
| 镜像不可变 | Digest 部署 | CD 模板 |
| 发布必须可追溯 | Commit / Build / Image Label | Admission / Audit |

---

## 安全

| 规范 | 自动化实现 | 强制点 |
|----|----|----|
| 未签名镜像禁止运行 | Kyverno verifyImages | Pod 创建失败 |
| 非官方 Registry 禁止 | Registry 白名单 | Admission |
| 必须生成 SBOM | syft 自动生成 | CI 阶段 |
| 严重漏洞禁止上线 | Trivy 扫描分级 | Pipeline Gate |
| 禁止特权容器 | Pod Security Policy / Pod Security Standards | Admission |

---

## 平台治理

| 规范 | 自动化 | 强制点 |
|----|----|----|
| 必须设置资源 limit | Kyverno Policy | Admission |
| 必须设置 owner label | Label 校验 | Admission |
| Namespace 隔离 | RBAC / Project | API Server |
| 生产环境审批 | Pipeline Approval | CD Gate |

---

## 发布质量 & SRE

| 规范 | 自动化 | 强制点 |
|----|----|----|
| 发布前检查 SLO | SLO API / Pipeline Gate | CD 阻断 |
| Error Budget 用尽禁止发布 | 自动计算 | CD 阻断 |
| 灰度失败自动回滚 | Argo Rollouts | 控制器 |
| 发布过程可观测 | Trace / Event | 平台内建 |

---

## 运行时安全

| 规范 | 自动化 | 强制点 |
|----|----|----|
| 容器必须 non-root | Pod Security / Kyverno | Admission |
| 异常行为可感知 | Falco / eBPF | 告警 |
| 运行行为可审计 | Audit Log | 集中存储 |

---

## 成本与资源

| 规范 | 自动化 | 强制点 |
|----|----|----|
| 资源必须可归属 | Label / Namespace | Admission |
| 防止资源浪费 | 利用率分析 | 告警 |
| 超预算提醒 | Budget Policy | 告警 |

---

## 自服务 & 开发者体验（DX）

| 规范 | 自动化 | 强制点 |
|----|----|----|
| 服务创建标准化 | Service Template / Scaffold | 平台入口 |
| CI / CD 自动开通 | Pipeline Scaffold | 平台 API |
| 监控自动接入 | 默认仪表盘 / Dashboard | 平台内建 |

---

## 最小落地集（MVP）

> 平台立刻可落地、最小争议、建立权威的规范集

1. 未签名镜像禁止运行  
2. 镜像必须来自平台 Registry  
3. 必须设置资源 limit  
4. 必须带 owner label  
5. Digest 部署，不允许 tag  

> 这五条规则是平台安全与治理的核心起点，能够快速建立「平台权威」。

---

## 附录

- **术语表**
  - **SBOM**：软件组成清单（Software Bill of Materials）
  - **SLO / SLI**：服务级目标与指标（Service Level Objective / Indicator）
  - **MVP**：最小可落地版本（Minimum Viable Product）

- **参考标准**
  - Kubernetes Security Best Practices
  - SLSA（Supply-chain Levels for Software Artifacts）
  - NIST DevSecOps Guidelines

---

> **核心思想**：  
> 平台工程的成熟标志，不是让大家更自由，而是让错误更难发生。  
> 规范 → 自动化 → 强制点 = 云原生平台真正价值所在。