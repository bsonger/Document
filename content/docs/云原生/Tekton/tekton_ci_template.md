
# Tekton CI Pipeline 模版（生产级）

> 最近更新：2025-05-09  
> 适用 **Kubernetes ≥ 1.26**、**Tekton Pipelines ≥ v1.0**  

该模版覆盖典型 CI 流程：**代码获取 → 代码扫描 → 镜像构建推送 → 镜像扫描**。  
所有 Task 与 Pipeline 以 YAML 形式提供，可直接 `kubectl apply -f` 部署，也可根据需要进行定制。

---

## 目录
1. 准备工作  
2. Secret 与 ServiceAccount  
3. Task 定义  
4. Pipeline 定义  
5. PipelineRun 示例  
6. 常见扩展（Chains 签名、集成 Triggers）  

---

## 1. 准备工作

```bash
# 创建工作命名空间
kubectl create ns ci

# 如果尚未安装 Tekton Pipelines
# 参考：https://tekton.dev/docs/pipelines/install/
```

---

## 2. Secret 与 ServiceAccount

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registry-auth   # 镜像仓库凭证
  namespace: ci
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: <BASE64>

---
apiVersion: v1
kind: Secret
metadata:
  name: sonar-token     # SonarQube Token
  namespace: ci
type: Opaque
stringData:
  token: <YOUR_SONAR_TOKEN>

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ci-builder
  namespace: ci
secrets:
  - name: registry-auth
  - name: sonar-token
```

> 🔒 **最佳实践**  
> - 通过 [External Secrets Operator](https://external-secrets.io/) 或 Vault 将凭证动态注入。  
> - 针对不同仓库/环境创建最小权限 ServiceAccount。

---

## 3. Task 定义

### 3.1 `git-clone`（引用官方 ClusterTask）

```yaml
# 官方 catalog: https://hub.tekton.dev/tekton/task/git-clone
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: git-clone
  namespace: ci
spec:
  params:
    - name: url
      type: string
    - name: revision
      type: string
      default: main
    - name: depth
      default: "1"
  workspaces:
    - name: output
  steps:
    - name: clone
      image: alpine/git:v2.44.0
      script: |
        git clone --depth=$(params.depth) -b $(params.revision) $(params.url) $(workspaces.output.path)
```

### 3.2 `sonar-code-scan`

```yaml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: sonar-code-scan
  namespace: ci
spec:
  params:
    - name: sonar-host-url
      type: string
    - name: project-key
      type: string
    - name: extra-args
      type: string
      default: ""
  workspaces:
    - name: source
  steps:
    - name: scan
      image: sonarsource/sonar-scanner-cli:latest
      env:
        - name: SONAR_TOKEN
          valueFrom:
            secretKeyRef:
              name: sonar-token
              key: token
      workingDir: $(workspaces.source.path)
      script: |
        sonar-scanner           -Dsonar.projectKey=$(params.project-key)           -Dsonar.sources=.           -Dsonar.host.url=$(params.sonar-host-url)           -Dsonar.login=$SONAR_TOKEN           $(params.extra-args)
```

### 3.3 `buildah-build-push`

```yaml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: buildah-build-push
  namespace: ci
spec:
  params:
    - name: image
      type: string     # e.g. harbor.mycorp.com/demo/web
    - name: tag
      type: string     # e.g. $(tt.params.revision) 或者 commit SHA
    - name: context
      type: string
      default: .
  workspaces:
    - name: source
  steps:
    - name: build-push
      image: quay.io/buildah/stable:v1.35
      securityContext:
        privileged: true          # 需要在安全策略允许
      env:
        - name: REGISTRY_AUTH
          valueFrom:
            secretKeyRef:
              name: registry-auth
              key: .dockerconfigjson
      workingDir: $(workspaces.source.path)
      script: |
        echo "$REGISTRY_AUTH" > /root/.docker/config.json
        IMAGE=$(params.image):$(params.tag)
        buildah bud --layers -t $IMAGE $(params.context)
        buildah push $IMAGE
```

### 3.4 `trivy-image-scan`

```yaml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: trivy-image-scan
  namespace: ci
spec:
  params:
    - name: image
      type: string
    - name: severity
      type: string
      default: HIGH,CRITICAL
    - name: exit-code
      type: string
      default: "1"
  steps:
    - name: scan
      image: aquasec/trivy:0.50.0
      script: |
        trivy image --exit-code $(params.exit-code) --severity $(params.severity) $(params.image)
```

---

## 4. Pipeline 定义

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: ci-pipeline
  namespace: ci
spec:
  workspaces:
    - name: shared-workspace
  params:
    - name: git-url
      type: string
    - name: git-revision
      type: string
      default: main
    - name: sonar-host-url
      type: string
    - name: sonar-project-key
      type: string
    - name: image-registry
      type: string         # harbor.mycorp.com/demo/web
    - name: image-tag
      type: string
      default: latest
  tasks:
    - name: clone
      taskRef:
        name: git-clone
      params:
        - name: url
          value: $(params.git-url)
        - name: revision
          value: $(params.git-revision)
      workspaces:
        - name: output
          workspace: shared-workspace

    - name: code-scan
      runAfter: ["clone"]
      taskRef:
        name: sonar-code-scan
      params:
        - name: sonar-host-url
          value: $(params.sonar-host-url)
        - name: project-key
          value: $(params.sonar-project-key)
      workspaces:
        - name: source
          workspace: shared-workspace

    - name: build-push
      runAfter: ["code-scan"]
      taskRef:
        name: buildah-build-push
      params:
        - name: image
          value: $(params.image-registry)
        - name: tag
          value: $(params.image-tag)
        - name: context
          value: .
      workspaces:
        - name: source
          workspace: shared-workspace

    - name: image-scan
      runAfter: ["build-push"]
      taskRef:
        name: trivy-image-scan
      params:
        - name: image
          value: $(params.image-registry):$(params.image-tag)
```

---

## 5. PipelineRun 示例

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  generateName: ci-run-
  namespace: ci
spec:
  serviceAccountName: ci-builder
  pipelineRef:
    name: ci-pipeline
  workspaces:
    - name: shared-workspace
      volumeClaimTemplate:
        metadata:
          name: ws-ci
        spec:
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: 1Gi
  params:
    - name: git-url
      value: https://github.com/mycorp/backend.git
    - name: git-revision
      value: main
    - name: sonar-host-url
      value: https://sonarqube.mycorp.com
    - name: sonar-project-key
      value: backend
    - name: image-registry
      value: harbor.mycorp.com/backend/backend-api
    - name: image-tag
      value: "$(gitrev)"
```

> ⚙️ `$(gitrev)` 示意：可通过在 TriggerBinding 中注入 commit SHA 创建动态 tag。  

---

## 6. 常见扩展

| 需求 | 做法 |
|------|------|
| **CI 完成后自动触发 CD** | 在 Triggers 上为 `PipelineRun` Success 设置 CloudEvent / ArgoCD Webhook |
| **镜像签名** | 安装 **Tekton Chains** (`tekton/chains`) 并启用 cosign |
| **并行测试** | 在 Pipeline 中新增 `unit-test`, `integration-test` Task，运行于 `runAfter: ["clone"]` 且互相并行 |
| **多架构镜像** | 使用 buildah/kaniko 的 `--platform linux/amd64,linux/arm64` 构建 |
| **资源隔离** | 基于 `workspaces` + PVC 或基于 `pvc.workspace.bindReadOnly` 限制 |

---

## 应用

```bash
kubectl apply -f tekton-ci-template.yaml          # 本文所有 YAML 合并保存后
# 或
kubectl apply -f task/
kubectl apply -f pipeline/ci-pipeline.yaml
kubectl create -f pipelinerun/ci-run.yaml
```

---

> 如需进一步自定义、接入企业内部制品库、支持缓存、或创建 Helm Chart，可联系作者。  

