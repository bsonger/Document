---
title: "introduce"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---

# cert-manager 介绍与原理文档

![cert-manager Logo](https://cert-manager.io/images/logo.svg)

## 一、什么是 cert-manager？

`cert-manager` 是一个 Kubernetes 插件，用于自动化管理和颁发 TLS 证书。它可以帮助用户从多个证书颁发机构（如 Let's Encrypt、HashiCorp Vault、自建 CA）自动请求和续签证书，从而确保集群中的服务通信安全。

官网地址：https://cert-manager.io/

## 二、作用

- 自动化 TLS 证书申请和续期
- 支持多种 Issuer（颁发者）类型，如 ACME（Let's Encrypt）、CA、Vault 等
- 统一 Kubernetes 中证书的管理方式
- 配合 Ingress 控制器使用，实现 HTTPS 自动化
- 与 Kubernetes 原生对象集成，如 Secret、Ingress、Service 等

## 三、核心概念

| 名称        | 说明                                                                 |
|-------------|----------------------------------------------------------------------|
| **Issuer**  | 表示一个证书颁发者，可为命名空间级别或全局（ClusterIssuer）         |
| **Certificate** | 表示一个证书请求，cert-manager 会据此去申请证书并生成 Secret     |
| **Challenge**   | ACME 协议中的校验挑战过程，用于验证域名所有权                   |
| **Order**       | 一个 ACME 协议中发起的证书申请订单                                |
| **Secret**      | 存储颁发后的证书和私钥，用于挂载到 Pod 或配置到 Ingress 中       |

## 四、工作原理

![cert-manager 流程图](https://cert-manager.io/docs/images/cert-manager-overview.svg)

## 五、架构图（Mermaid）

```mermaid
graph TD
  A[User / YAML] -->|Create Certificate| B[cert-manager Controller]
  B --> C[Check Issuer or ClusterIssuer]
  C --> D{Issuer Type}
  D -->|ACME| E[ACME Solver]
  E --> F[Challenge + Order]
  F --> G[ACME CA (e.g., Let's Encrypt)]
  D -->|CA / Vault| H[Internal CA / Vault API]
  G --> I[Signed Certificate]
  H --> I
  I --> J[Secret Created / Updated]
  J --> K[Used by Ingress / Pod]
```

## 六、常见使用场景

- 配合 NGINX Ingress Controller 实现自动 HTTPS
- 自动颁发和续期 Pod 之间 mTLS 所需的证书
- 对接内部 CA，提供可信任的内部证书签发服务
- 与 Vault 集成，实现高度安全的证书存储和访问

## 七、参考资源

- [cert-manager 官方文档](https://cert-manager.io/docs/)
- [Kubernetes Ingress TLS 配置指南](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)
