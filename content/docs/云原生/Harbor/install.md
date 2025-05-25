# 🚀 Harbor 生产环境安装与 Keycloak 对接指南

## 🏗️ 一、生产环境 Harbor 安装步骤

### 1. 系统准备
- 操作系统：Ubuntu 20.04 / CentOS 7+
- 环境要求：
  - CPU ≥ 2 核
  - 内存 ≥ 4 GB
  - 磁盘 ≥ 40 GB
- 安装依赖：
  ```bash
  sudo apt update && sudo apt install docker.io docker-compose -y
  ```

### 2. 下载 Harbor 安装包
```bash
wget https://github.com/goharbor/harbor/releases/download/v2.10.0/harbor-offline-installer-v2.10.0.tgz

# 解压
 tar -zxvf harbor-offline-installer-*.tgz
cd harbor
```

### 3. 配置 harbor.yml
编辑 `harbor.yml` 关键配置项：
```yaml
hostname: harbor.yourdomain.com

https:
  port: 443
  certificate: /path/to/cert.pem
  private_key: /path/to/key.pem

harbor_admin_password: YourStrongAdminPass

auth_mode: oidc
oidc_provider_name: keycloak
oidc_endpoint: https://keycloak.yourdomain.com/realms/harbor/.well-known/openid-configuration
oidc_client_id: harbor-client
oidc_client_secret: <your-client-secret>
oidc_scope: openid,profile,email,groups
oidc_verify_cert: true
oidc_auto_onboard: true
```

### 4. 安装 Harbor
```bash
./install.sh
```

---

## 🧩 二、配置 Keycloak 与 Harbor 对接

### 1. 创建 Realm
- 名称建议为 `harbor`

### 2. 创建 Client
- Client ID: `harbor-client`
- Access Type: `confidential`
- Root URL: `https://harbor.yourdomain.com`
- Redirect URI: `https://harbor.yourdomain.com/c/oidc/callback`

### 3. 启用功能
- 启用以下选项：
  - Standard Flow Enabled
  - Direct Access Grants Enabled
  - Service Accounts Enabled

### 4. 映射用户组（可选）
- 添加 Mapper：
  - Name: `groups`
  - Token Claim Name: `groups`
  - Claim JSON Type: `String`
  - Add to ID token / access token: ✅

---

## 🔐 三、安全建议

| 项目               | 建议                                               |
|--------------------|----------------------------------------------------|
| TLS 证书           | 使用 CA 签发证书                                   |
| 镜像扫描           | 启用 Trivy 漏洞扫描                                |
| 数据库外置         | 使用外部 PostgreSQL 提升可维护性                  |
| 镜像签名           | 配合 Cosign / Notary 实现镜像签名                  |
| RBAC 权限控制      | Keycloak 分组 + Harbor 项目权限                    |
| 日志与审计         | 启用 `audit.log` 并定期备份                        |

---

## ✅ 四、验证流程
1. 启动 Harbor 并访问 `https://harbor.yourdomain.com`
2. 点击登录，将跳转到 Keycloak 页面
3. 登录后自动完成注册并进入 Harbor 控制台

---

## 📦 五、补充建议
- Harbor 支持高可用部署，可配合 Nginx Ingress + Keepalived 实现
- Harbor 支持 Helm Charts 管理，可结合 GitOps 工具使用
- 建议定期更新 Trivy 数据库，保持漏洞扫描能力的时效性
