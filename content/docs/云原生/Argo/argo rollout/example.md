
# 🚀 Argo CD 生产环境最复杂示例（多服务、多集群、多环境、Helm + Rollouts）

本示例展示如何使用 Argo CD 管理多个服务、环境和集群，结合 Helm、Rollouts、Project、RBAC、ApplicationSet 构建完整 GitOps 流水线。

---

## 📁 Git 仓库结构约定

```
apps/
├── svc-a/
│   └── charts/
│       ├── dev/values.yaml
│       ├── stage/values.yaml
│       └── prod/values.yaml
├── svc-b/
│   └── charts/
│       ├── dev/values.yaml
│       └── prod/values.yaml
└── rollout-templates/
    └── nginx-rollout.yaml
```

---

## 🧩 1. AppProject 定义（多服务隔离）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: microservices
  namespace: argocd
spec:
  description: Production microservice apps
  sourceRepos:
    - https://github.com/your-org/gitops-repo.git
  destinations:
    - namespace: svc-*
      server: https://kubernetes.default.svc
    - namespace: svc-*
      server: https://k8s-prod.example.com
  clusterResourceWhitelist:
    - group: "*"
      kind: "*"
  namespaceResourceWhitelist:
    - group: "*"
      kind: "*"
```

---

## ⚙️ 2. ApplicationSet（服务 x 环境 x 集群 x Rollout）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: multi-svc-prod
  namespace: argocd
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/your-org/gitops-repo.git
              revision: main
              directories:
                - path: apps/*/charts/*
          - list:
              elements:
                - env: dev
                  cluster: https://kubernetes.default.svc
                - env: prod
                  cluster: https://k8s-prod.example.com
  template:
    metadata:
      name: '{{path[1]}}-{{env}}'
    spec:
      project: microservices
      source:
        repoURL: https://github.com/your-org/gitops-repo.git
        targetRevision: main
        path: apps/{{path[1]}}/charts/{{env}}
        helm:
          valueFiles:
            - values.yaml
      destination:
        server: '{{cluster}}'
        namespace: svc-{{path[1]}}-{{env}}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
          - ApplyOutOfSyncOnly=true
```

---

## 🔁 3. 示例服务使用 Rollouts（Helm Chart 中包含 rollout.yaml）

在 `values.yaml` 中控制是否启用 Rollout 类型资源替代 Deployment。

```yaml
rollout:
  enabled: true
  strategy: canary
  canary:
    steps:
      - setWeight: 10
      - pause: {duration: 1m}
      - setWeight: 100
    analysis:
      enabled: true
      templateName: success-rate-check
```

Chart 模板结构参考 `nginx-rollout.yaml`：

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: {{ .Release.Name }}
spec:
  strategy:
    canary:
      steps:
        - setWeight: 10
        - pause: {duration: 1m}
```

---

## 🔐 4. RBAC 限制

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.csv: |
    g, devops-team, role:microservices-admin
    p, role:microservices-admin, applications, *, microservices/*, allow
```

---

## 🔍 5. 集成建议

- Prometheus + Grafana：用于监控 Rollout 状态与指标分析
- Istio：用于金丝雀流量控制
- Dex + OIDC：用户认证 + 单点登录
- Argo Notifications：通知部署状态变化
- Loki/EFK：日志审计
- Cert-Manager + Ingress：HTTPS 访问 Argo CD 界面

---

## 🧠 总结功能覆盖

| 模块 | 是否包含 |
|------|----------|
| 多服务支持 | ✅ |
| 多环境（dev/prod） | ✅ |
| 多集群部署 | ✅ |
| Helm 管理配置 | ✅ |
| Canary Rollout + Prometheus | ✅ |
| 自动同步、自愈、命名空间自动创建 | ✅ |
| RBAC + Project + SSO | ✅ |
| 可观测性与通知 | ✅ |
| ApplicationSet 动态生成 | ✅ |

---

## 📚 参考文档

- Argo CD ApplicationSet：https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/
- Argo Rollouts：https://argo-rollouts.readthedocs.io/
- Prometheus Metrics + AnalysisRun：https://argo-rollouts.readthedocs.io/en/stable/features/analysis/
