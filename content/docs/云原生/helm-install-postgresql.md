
# 使用 Helm 安装 PostgreSQL 到 Kubernetes

## 前提条件

- 已安装并配置好的 Kubernetes 集群
- 本地已安装并配置好 `kubectl`
- 安装了 Helm v3+

## 步骤一：添加 Bitnami Helm 仓库

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

## 步骤二：创建命名空间（可选）

```bash
kubectl create namespace postgres
```

## 步骤三：安装 PostgreSQL

```bash
helm install my-postgresql \
  --namespace postgres \
  --set auth.postgresPassword=your-password \
  bitnami/postgresql
```

你可以根据需要修改参数，如：

- `auth.postgresPassword`：设置 PostgreSQL 的密码
- `primary.persistence.enabled=true/false`：是否启用数据持久化
- `primary.resources.requests`：设置 CPU/内存请求资源

## 步骤四：验证安装

```bash
kubectl get pods -n postgres
kubectl get svc -n postgres
```

## 步骤五：连接到 PostgreSQL

### 方式一：通过 Pod 内部连接

```bash
kubectl run -i --tty --rm psql-client --image=bitnami/postgresql --namespace postgres \
  --env="PGPASSWORD=your-password" \
  --command -- psql -h my-postgresql.postgres.svc.cluster.local -U postgres
```

### 方式二：通过端口转发（本地连接）

```bash
kubectl port-forward --namespace postgres svc/my-postgresql 5432:5432
```

然后使用本地工具连接：

```bash
psql -h 127.0.0.1 -U postgres -p 5432
```

密码为 `your-password`

## 步骤六：卸载 PostgreSQL（可选）

```bash
helm uninstall my-postgresql -n postgres
kubectl delete namespace postgres
```

## 附录：常用 Helm 参数

| 参数 | 描述 | 示例 |
|------|------|------|
| `auth.postgresPassword` | 设置 PostgreSQL 密码 | `your-password` |
| `primary.persistence.enabled` | 启用持久化 | `true` |
| `primary.persistence.size` | 存储大小 | `8Gi` |
| `primary.resources.requests.memory` | 内存请求 | `512Mi` |
| `primary.resources.requests.cpu` | CPU 请求 | `250m` |

更多参数详见：[Bitnami PostgreSQL Helm Chart 文档](https://artifacthub.io/packages/helm/bitnami/postgresql)
