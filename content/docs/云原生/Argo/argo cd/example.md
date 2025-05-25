---
title: "example"
weight: 4
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---
# 🚀 Argo CD 多集群、多环境、Helm 管理 + Project 隔离示例（生产级复杂实践）

本示例展示如何使用 Argo CD + ApplicationSet 管理多个集群、多环境、多服务，结合 Helm、Argo CD Project 和 RBAC，构建生产级 GitOps 平台。

---

## 📁 Git 仓库结构假设

```
git-repo/
└── services/
    ├── svc-a/
    │   └── charts/
    │       ├── dev/values.yaml
    │       ├── stage/values.yaml
    │       └── prod/values.yaml
    ├── svc-b/
    │   └── charts/
    │       ├── dev/values.yaml
    │       └── prod/values.yaml
```

---

## 🧩 1. 创建 Argo CD Project（隔离服务与环境）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: svc-a
  namespace: argocd
spec:
  description: svc-a service deployments
  sourceRepos:
    - https://github.com/your-org/your-repo.git
  destinations:
    - namespace: svc-a-*
      server: https://kubernetes.default.svc
    - namespace: svc-a-*
      server: https://k8s-stage.example.com
    - namespace: svc-a-*
      server: https://k8s-prod.example.com
  clusterResourceWhitelist:
    - group: "*"
      kind: "*"
  namespaceResourceWhitelist:
    - group: "*"
      kind: "*"
```

---

## ⚙️ 2. ApplicationSet 示例（支持多环境 + 多集群 + Helm + Project）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: svc-a-apps
  namespace: argocd
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/your-org/your-repo.git
              revision: main
              directories:
                - path: services/svc-a/charts/*
          - list:
              elements:
                - env: dev
                  cluster: https://kubernetes.default.svc
                - env: stage
                  cluster: https://k8s-stage.example.com
                - env: prod
                  cluster: https://k8s-prod.example.com
  template:
    metadata:
      name: svc-a-{{env}}
    spec:
      project: svc-a
      source:
        repoURL: https://github.com/your-org/your-repo.git
        targetRevision: main
        path: services/svc-a/charts/{{env}}
        helm:
          valueFiles:
            - values.yaml
      destination:
        server: '{{cluster}}'
        namespace: svc-a-{{env}}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
```

---

## 🔒 3. RBAC 示例（可选）

为特定用户组限制访问 `svc-a` 项目与环境：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.csv: |
    g, dev-team, role:svc-a-dev
    p, role:svc-a-dev, applications, get, svc-a-*, allow
    p, role:svc-a-dev, projects, get, svc-a, allow
```

---

## ✅ 高级功能回顾

| 功能 | 描述 |
|------|------|
| Project 资源隔离 | 每个服务独立 Project，可设置独立权限 |
| 多集群部署 | 根据环境匹配不同 Kubernetes API server |
| 多环境支持 | 使用 Helm 的子目录 values 管理 dev/stage/prod |
| Matrix Generator | 支持服务 + 环境组合生成多个应用 |
| 自动同步与自愈 | `automated + selfHeal` 提升可靠性 |
| 命名空间自动创建 | 支持 `CreateNamespace=true` |
| RBAC 控制 | 精细化控制用户访问权限 |

---

## 📦 生产部署注意事项

- 集群需提前通过 `argocd cluster add` 添加并设置 RBAC
- 每个服务建议用独立 Project 管理
- Helm charts 存放路径应一致规范
- 使用专有 Git 分支管理环境（可扩展 Git Generator 条件）
- 结合通知系统监控同步失败状态

---

## 🔗 推荐资源

- Project 配置文档：https://argo-cd.readthedocs.io/en/stable/user-guide/projects/
- Matrix Generator 文档：https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/MATRIX-Generator/
- Helm 支持说明：https://argo-cd.readthedocs.io/en/stable/user-guide/helm/
