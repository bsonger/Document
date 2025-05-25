---
title: "e2e_flow_full"
weight: 6
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false

---

# 🌐 Nginx 全家桶 GitOps – Dev & Prod 金丝雀发布

利用 **Argo Events → Argo Workflows → Harbor → Argo CD → Argo Rollouts**  
实现 *代码 Push → Dev Canary → 审批 → Prod Canary* 的端到端流水线，  
并通过 **Prometheus Analysis** 与两级人工审批保障生产安全。
---

## 🖼️ 架构图

```mermaid
flowchart TD
    P[Git Push] --> ES(EventSource<br/>Git Webhook)
    ES --> EB[EventBus NATS (3)]
    EB --> S(Sensor)
    S --> WF[CI Workflow<br/>Kaniko + Helm + 两级审批]
    WF --> IMG[(Harbor 镜像)]
    WF --> CH[(Harbor Chart)]
    WF -.sync-dev.-> CD1[Argo CD Dev<br/>Rollout]
    WF -.sync-prod.-> CD2[Argo CD Prod<br/>Rollout]
    CD1 -->|25→100 %| VS1(Istio VS Dev)
    CD2 -->|10→50→100 %| VS2(Istio VS Prod)
    Prom[Prometheus] -.metrics.-> VS1 & VS2
```

## 📂 目录结构
```text
nginx-gitops/
├── Dockerfile
├── chart/
│   ├── Chart.yaml
│   ├── values.yaml          # 公共
│   ├── values-dev.yaml      # Dev
│   ├── values-prod.yaml     # Prod
│   └── templates/
│       ├── rollout.yaml
│       └── service.yaml
└── argo/
    ├── eventbus.yaml
    ├── eventsource-git.yaml
    ├── sensor-ci.yaml
    ├── workflowtemplate-ci-harbor.yaml
    ├── analysis/success-rate.yaml
    └── applications/
        ├── nginx-dev.yaml
        └── nginx-prod.yaml
```

## 1️⃣ EventBus（NATS）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventBus
metadata:
  name: default
  namespace: argo-events
spec:
  nats:
    native:
      replicas: 3
      auth: token
```

---

## 2️⃣ EventSource（Git Webhook）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventSource
metadata:
  name: git-webhook
  namespace: argo-events
spec:
  service:
    ports:
      - port: 12000
        targetPort: 12000
  github:
    git-push:
      endpoint: /push
      port: 12000
      insecure: false
      secretToken:
        name: github-secret
        key: token
```

---

## 3️⃣ Sensor（调用 WorkflowTemplate）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: ci-sensor
  namespace: argo-events
spec:
  eventBusName: default
  dependencies:
    - name: push
      eventSourceName: git-webhook
      eventName: git-push
  triggers:
    - template:
        name: ci-workflow
        k8s:
          group: argoproj.io
          version: v1alpha1
          resource: workflows
          operation: create
          source:
            resource:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: ci-pipeline-
              spec:
                workflowTemplateRef:
                  name: ci-template
```

---

## 4️⃣ WorkflowTemplate（CI → 审批 → GitOps）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: ci-nginx-harbor
  namespace: argo
spec:
  entrypoint: pipeline
  serviceAccountName: argo-ci
  templates:
    - name: pipeline
      dag:
        tasks:
          - name: build
            template: kaniko
          - name: package
            dependencies: [build]
            template: helm-package
          - name: push-chart
            dependencies: [package]
            template: helm-push
          - name: approve-dev
            dependencies: [push-chart]
            template: pause
          - name: sync-dev
            dependencies: [approve-dev]
            template: argocd-sync
            arguments:
              parameters: [{name: app, value: nginx-dev}]
          - name: approve-prod
            dependencies: [sync-dev]
            template: pause
          - name: sync-prod
            dependencies: [approve-prod]
            template: argocd-sync
            arguments:
              parameters: [{name: app, value: nginx-prod}]

    # -- Kaniko 构建并推镜像 --
    - name: kaniko
      outputs:
        parameters:
          - name: tag
            valueFrom:
              path: /tmp/tag
      container:
        image: gcr.io/kaniko-project/executor:v1.20.0
        env:
          - name: TAG
            value: "{{workflow.labels.git-sha}}"
        args:
          - --dockerfile=Dockerfile
          - --destination=harbor.example.com/library/nginx:${TAG}
        lifecycle:
          postStart:
            exec:
              command: ["sh","-c","echo ${TAG} >/tmp/tag"]
        volumeMounts:
          - name: harbor-auth
            mountPath: /kaniko/.docker
      volumes:
        - name: harbor-auth
          secret:
            secretName: harbor-secret

    # -- Helm Package（写 Dev/Prod） --
    - name: helm-package
      inputs:
        parameters:
          - name: tag
            value: "{{tasks.build.outputs.parameters.tag}}"
      script:
        image: alpine/helm:3.13.0
        command: [sh]
        source: |
          TAG=$0
          cp -R /workspace/chart /tmp/chart && cd /tmp/chart
          for env in dev prod; do
            yq -i '.image.tag = strenv(TAG)' values-${env}.yaml
          done
          yq -i '.version=strenv(TAG) | .appVersion=strenv(TAG)' Chart.yaml
          helm dependency update .
          helm package . -d /tmp/out
          echo /tmp/out/$(ls /tmp/out) >/tmp/chart
        env:
          - name: 0
            value: "{{inputs.parameters.tag}}"
      outputs:
        artifacts:
          - name: chart
            path: /tmp/chart

    # -- Helm Push 到 Harbor OCI --
    - name: helm-push
      inputs:
        artifacts:
          - name: chart
            from: "{{tasks.helm-package.outputs.artifacts.chart}}"
      script:
        image: alpine/helm:3.13.0
        command: [sh]
        source: |
          helm registry login harbor.example.com -u $USER -p $PASS
          helm push $(cat {{inputs.artifacts.chart}}) oci://harbor.example.com/chartrepo
        env:
          - name: USER
            value: robot$library+ci
          - name: PASS
            valueFrom:
              secretKeyRef:
                name: harbor-secret
                key: .dockerconfigjson

    - name: pause
      suspend: {}

    - name: argocd-sync
      inputs:
        parameters:
          - name: app
      container:
        image: argoproj/argocd:v2.10.0
        command:
          - argocd
          - app
          - sync
          - "{{inputs.parameters.app}}"
          - --grpc-web
          - --insecure
        env:
          - name: ARGOCD_AUTH_TOKEN
            valueFrom:
              secretKeyRef:
                name: argocd-token
                key: token
```

---

## 5️⃣ Argo CD Application + Rollout（金丝雀）

`values.yaml` 里使用 Rollout，示例略（见前文 Rollouts 示例）。  
将 `autoPromotionEnabled: false` 可以让运维手动切流量。  

---

### 🔑 关键点

| 步骤 | 技术 | 备注 |
|------|------|------|
| 事件捕获 | **Argo Events** | Git Webhook |
| CI/Image | **Argo Workflows + Kaniko** | 无 docker-in-docker |
| 审批 | **Workflow Suspend** | UI / CLI Resume |
| GitOps | **Workflow Git commit** | 改 chart tag |
| CD & Canary | **Argo CD + Rollouts** | Dev / Prod 集群 |
| 指标分析 | **Prometheus** | Rollouts Analysis |

---

> 按需扩展更多服务可用 **ApplicationSet Matrix**。  
> 如需完整示例仓库，请将上面 YAML 放入新版 GitOps 存储结构。
