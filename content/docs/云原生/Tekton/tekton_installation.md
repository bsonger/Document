# Tekton 安装全攻略（Pipelines v1.0.0 + Triggers v0.31.0 + Dashboard v0.56.0）

> 最近更新：2025-05-09  
> 适用 Kubernetes ≥ 1.26

---

## 目录
1. [前提条件](#前提条件)  
2. [环境准备](#环境准备)  
3. [安装 Tekton Pipelines](#安装-tekton-pipelines)  
4. [安装 Tekton Triggers](#安装-tekton-triggers)  
5. [安装 Tekton Dashboard](#安装-tekton-dashboard)  
6. [安装 Tekton CLI (`tkn`)](#安装-tekton-cli-tkn)  
7. [RBAC 与安全加固](#rbac-与安全加固)  
8. [高可用 (HA) 部署选项](#高可用-ha-部署选项)  
9. [验证安装](#验证安装)  
10. [升级 / 卸载](#升级--卸载)  
11. [常见问题 (FAQ)](#常见问题-faq)  

---

## 前提条件

| 组件 | 版本 | 说明 |
|------|------|------|
| Kubernetes | ≥ 1.26 | 已开启 `MutatingAdmissionWebhook`、`ValidatingAdmissionWebhook` |
| `kubectl`  | ≥ 1.26 | 集群管理员权限 |
| Ingress Controller | 任意 | 用于暴露 Dashboard，可选 |
| StorageClass | 默认 | `PipelineRun` 可能需要持久化工作目录 |

---

## 环境准备

```bash
# 创建专用命名空间
kubectl create ns tekton-pipelines

# 建议给节点设置 registry-mirror & image-pull-secret 以加速拉取镜像
```

---

## 安装 Tekton Pipelines

### 获取最新清单

```bash
# 安装最新正式版（v1.0.0）
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```

> **Tips**  
> - 若使用私有镜像仓库，请先 `kubectl set image` 或编辑 YAML 替换镜像前缀  
> - 如需 **仅 CRD 安装**：`release.notags.yaml`  
> - 安装位置默认为 `tekton-pipelines` 命名空间

### 等待 Pod 就绪

```bash
watch kubectl get pods -n tekton-pipelines
```

---

## 安装 Tekton Triggers

```bash
# 安装 Triggers v0.31.0
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml
```

### 校验

```bash
kubectl get pods -n tekton-pipelines -l app.kubernetes.io/part-of=triggers
```

---

## 安装 Tekton Dashboard

Dashboard 支持两种方式：

### 方式 A：官方安装脚本（推荐，支持自定义）

```bash
curl -sL https://raw.githubusercontent.com/tektoncd/dashboard/main/scripts/release-installer \
  | bash -s -- install v0.56.0 --namespace tekton-pipelines --read-write
```

常用参数：

| 参数 | 说明 |
|------|------|
| `install <ver>` | 指定版本或 `latest` |
| `--namespace`   | 安装目标命名空间 |
| `--read-only` / `--read-write` | 访问模式 |
| `--ingress-class <class>` | 自动创建 Ingress |

### 方式 B：YAML

```bash
kubectl apply -f https://storage.googleapis.com/tekton-releases/dashboard/latest/release.yaml
```

#### 访问 Dashboard

```bash
# 暴露 NodePort (示例)
kubectl patch svc tekton-dashboard -n tekton-pipelines -p '{{"spec": {{"type": "NodePort"}}}}'
minikube service tekton-dashboard -n tekton-pipelines
```

或创建 Ingress / LoadBalancer 根据集群情况选择。

---

## 安装 Tekton CLI (`tkn`)

Linux/macOS (x86_64 & arm64):

```bash
curl -LO https://github.com/tektoncd/cli/releases/latest/download/tkn_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.tar.gz
sudo tar -C /usr/local/bin -zxvf tkn_*_amd64.tar.gz tkn
```

验证：

```bash
tkn version
```

---

## RBAC 与安全加固

```bash
# 仅允许特定命名空间触发 PipelineRun
kubectl create clusterrole tekton-triggers-min \
  --verb=get,list,watch,create,update,patch,delete \
  --resource=eventlisteners,triggers
kubectl create clusterrolebinding tekton-triggers-min \
  --clusterrole=tekton-triggers-min --serviceaccount=tekton-pipelines:tekton-triggers-controller
```

- 对接 **OIDC**：在 Dashboard Ingress 上配置 `oauth2-proxy`。  
- 使用 `serviceAccountName` 绑定最小权限的 `TaskRun`s。  
- 若集群启用了 PSP/PodSecurity 标准，需为 `tekton-pipelines` 命名空间打 `pod-security.kubernetes.io/enforce=privileged` 或自定义策略。

---

## 高可用 (HA) 部署选项

| 组件 | 建议副本数 | 额外说明 |
|------|-----------|----------|
| `tekton-pipelines-controller` | 2+ | 需开启 leader election |
| `tekton-triggers-controller` | 2+ | 建议启用 leader election |
| `tekton-dashboard` | 2+ | 无状态，可加 HPA |
| `webhook` | 2+ | Pipelines & Triggers Webhook |

> 使用 K8s 原生 **HPA** 结合自定义 Metrics 或 **Keda** 弹性伸缩。

---

## 验证安装

1. 创建示例 Task & Pipeline：

```bash
kubectl apply -f https://raw.githubusercontent.com/tektoncd/pipeline/main/docs/samples/hello-world/hello-world.yaml
tkn pipeline start hello-world-pipeline --param WHO=tekton --showlog
```

2. 创建示例 Trigger（GitHub Push）：

```bash
kubectl apply -f https://raw.githubusercontent.com/tektoncd/triggers/main/docs/examples/v1beta1/github/README.md
```

触发 `push` 后可在 Dashboard 查看 PipelineRun 进度。

---

## 升级 / 卸载

### 升级

```bash
# 先备份 CR 数据
kubectl get pipelineruns,taskruns,triggertemplates -A -o yaml > backup.yaml

# Pipelines
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v<old>/release.yaml --prune -l tekton.dev/release=previous
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
# Triggers 同理
```

> **注意**：确保目标版本与 Dashboard/Triggers 兼容。

### 卸载

```bash
# CRD 与命名空间全部删除
kubectl delete ns tekton-pipelines
kubectl delete crd $(kubectl get crd -o name | grep tekton.dev)
```

---

## 常见问题 (FAQ)

| 问题 | 解决方案 |
|------|----------|
| Pod Pending, 报 “no space left on device” | 检查 PVC 与节点磁盘空间 |
| Dashboard 403 Forbidden | 配置正确的 `ClusterRole` 与 Ingress OIDC |
| GitHub webhook 无法被触发 | 检查 `EventListener` Service 外部可达性，确认 `Secret` 是否包含正确 token |
| 镜像拉取超时 | 使用私有镜像仓库或镜像加速器 |

---

### 参考链接

- Tekton 官方文档：<https://tekton.dev/docs/>  
- Pipelines 安装文档：<https://tekton.dev/docs/pipelines/install/>  
- Triggers 安装文档：<https://tekton.dev/docs/triggers/install/>  
- Dashboard 安装文档：<https://tekton.dev/docs/dashboard/install/>  

---

> 如需进阶示例（如 GitOps、KubeVirt CI、内置 Chains 签名验证）或企业级 Helm Chart，请联系维护者。  
