
# 🚀 Argo Events 支持的 Trigger 类型 + Sensor + EventBus 示例合集

---

## ✅ EventBus（NATS 默认）示例

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

## ✅ Sensor + Trigger 示例合集

---

### 📌 1. Kubernetes Trigger - 创建 Workflow

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: workflow-sensor
  namespace: argo-events
spec:
  dependencies:
    - name: webhook-dep
      eventSourceName: webhook-source
      eventName: webhook
  triggers:
    - template:
        name: create-workflow
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
                generateName: example-workflow-
              spec:
                entrypoint: whalesay
                templates:
                  - name: whalesay
                    container:
                      image: docker/whalesay
                      command: [cowsay]
                      args: ["Hello from Argo Events"]
```

---

### 📌 2. HTTP Trigger - 调用外部 Webhook

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: http-trigger-sensor
spec:
  dependencies:
    - name: webhook-dep
      eventSourceName: webhook-source
      eventName: webhook
  triggers:
    - template:
        name: call-http
        http:
          url: https://hooks.example.com/build
          method: POST
          payload:
            - src:
                dependencyName: webhook-dep
                dataKey: body
              dest: body.data
          headers:
            - name: Content-Type
              value: application/json
```

---

### 📌 3. Kafka Trigger - 发送消息

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: kafka-trigger-sensor
spec:
  dependencies:
    - name: kafka-dep
      eventSourceName: kafka-source
      eventName: my-topic
  triggers:
    - template:
        name: kafka-trigger
        kafka:
          url: kafka.default.svc.cluster.local:9092
          topic: result-topic
          partition: 0
          parameters:
            - src:
                dependencyName: kafka-dep
                dataKey: body.message
              dest: payload
```

---

### 📌 4. AWS Lambda Trigger - 执行函数

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: lambda-trigger-sensor
spec:
  dependencies:
    - name: webhook-dep
      eventSourceName: webhook-source
      eventName: webhook
  triggers:
    - template:
        name: call-lambda
        awsLambda:
          functionName: myLambdaFunction
          region: us-west-2
          accessKey:
            key: accesskey
            name: aws-creds
          secretKey:
            key: secretkey
            name: aws-creds
          payload:
            - src:
                dependencyName: webhook-dep
                dataKey: body
              dest: body.message
```

---

### 📌 5. Argo Workflow Trigger（简写方式）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: simple-workflow-sensor
spec:
  dependencies:
    - name: git-push
      eventSourceName: github-source
      eventName: push
  triggers:
    - template:
        name: simple-workflow
        argoWorkflow:
          group: argoproj.io
          version: v1alpha1
          resource: workflows
          operation: create
          source:
            inline:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: gitops-
              spec:
                entrypoint: hello
                templates:
                  - name: hello
                    container:
                      image: alpine
                      command: ["echo"]
                      args: ["Git pushed, deploying..."]
```

---

## 📚 总结

每种 Trigger 类型都可配合任意 EventSource，如 Webhook/Kafka/S3/Cron 等使用。EventBus 提供传输通道（默认 NATS），Sensor 管理依赖和事件逻辑，Trigger 执行实际动作。

更多模板参考：https://github.com/argoproj/argo-events/tree/master/examples
