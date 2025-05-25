# 🧪 Harbor 最复杂的生产使用示例

本示例演示在企业级生产环境中使用 Harbor 的复杂集成与最佳实践。该案例包括多数据中心部署、Keycloak 单点登录、CI/CD、镜像签名、漏洞扫描、多租户、异地灾备等完整场景。

## 🏭 场景背景
- 企业规模：跨国公司
- 数据中心：北京、上海、新加坡
- 用户管理：Keycloak + AD 同步
- 镜像类型：容器镜像、Helm Chart、AI 模型
- DevOps 工具：GitLab CI、Argo CD、Cosign
- 安全要求：镜像签名、漏洞扫描、RBAC、审计日志

---

## 🧰 组件部署架构
```
+---------------------------+         +---------------------------+
|     Developer (GitLab)   |         |   SSO (Keycloak + AD)    |
+-------------+-------------+         +-------------+-------------+
              |                                     |
              v                                     |
     GitLab CI/CD Runner                            |
              |                                     |
              v                                     v
     +--------+---------+               +-----------+-----------+
     |    Harbor (主)   | <------------> |    Keycloak Realm    |
     |  harbor.beijing  |               |     realm: harbor     |
     +--------+---------+               +------------------------+
              |
        镜像复制策略 (双向)
              v
     +--------+---------+
     |  Harbor (备)     |
     | harbor.singapore |
     +------------------+
```

---

## 🔐 安全配置

### 1. OIDC 登录与权限控制
- 启用 Keycloak OIDC 模式
- Harbor 设置为 `auth_mode: oidc`
- Keycloak 同步 AD 用户，分组映射 Harbor 项目
- Harbor 内部使用 RBAC 精细化授权（例如 `project-admin`, `developer`, `guest`）

### 2. 镜像签名（Cosign）
- CI/CD 阶段使用 Cosign 对构建好的镜像签名：
  ```bash
  cosign sign --key k8s://default/mykey docker.io/acme/app:1.0.0
  ```
- Harbor 启用镜像验证策略，确保部署前必须验证签名

### 3. 漏洞扫描（Trivy）
- Harbor 启用 Trivy Scanner
- 配置上传即扫描策略
- 结合 Argo CD 实现漏洞阻断部署

### 4. 审计日志与存储加密
- 审计日志开启，导出至 ELK 系统
- Harbor 存储后端为 S3，开启服务端加密（SSE）

---

## 🔄 镜像同步策略

- Harbor 主节点（北京）配置多个 replication rule：
  - 同步至上海 Harbor（延迟 5s）
  - 同步至新加坡 Harbor（用于灾备）
- 同步模式：Push-based，定时触发
- 支持双向同步（防止单点）

---

## 🚀 DevOps 流程集成

### CI 阶段（GitLab）
```yaml
build:
  script:
    - docker build -t harbor.beijing.company.com/dev/app:$CI_COMMIT_SHA .
    - docker push harbor.beijing.company.com/dev/app:$CI_COMMIT_SHA
    - cosign sign --key k8s://default/signkey harbor.beijing.company.com/dev/app:$CI_COMMIT_SHA
```

### CD 阶段（Argo CD）
- Argo CD 部署前校验镜像签名和扫描结果：
  - 使用 Gatekeeper + OPA 策略
  - 不合规镜像拒绝部署

---

## 📦 多类型镜像支持
- 容器镜像（标准）
- Helm Chart 仓库（开启 chartmuseum）
- AI 模型仓库（通过 OCI artifact 支持自定义类型）

---

## 📊 多租户管理
- 每个 BU（业务单元）对应一个 Harbor 项目
- 权限通过 Keycloak Group 自动同步
- 各项目启用资源配额（Quota）控制

---

## 🌐 灾备与高可用
- Harbor 主节点高可用部署（3 节点 Nginx + Keepalived + shared DB）
- 镜像异地同步 + 数据定期备份
- Redis、数据库外部化，便于迁移和监控

---

## ✅ 总结
此示例演示了如何在真实、复杂的企业生产环境中部署 Harbor，并与身份管理、CI/CD、安全机制、镜像治理等体系深度集成，构建一个稳定、安全、高效的镜像管理平台。
