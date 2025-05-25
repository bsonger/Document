---
title: "introduce"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---

# 📘 Argo CD Introduce

## 🧩 Argo CD 架构图与组件

![Argo CD 架构图](https://argo-cd.readthedocs.io/en/stable/assets/argo-cd-architecture.png)

### 核心组件解释：

#### API Server 🖥️

The API server is a GRPC/REST server which exposes the API consumed by the Web UI, CLI, and CI/CD systems. It has the following responsibilities:

1. Application management and status reporting
2. Invoking of application operations (e.g., sync, rollback, user‑defined actions)
3. Repository and cluster credential management (stored as Kubernetes Secrets)
4. Authentication and auth delegation to external identity providers
5. RBAC enforcement
6. Listener/forwarder for Git webhook events

---

#### Repository Server 📦

The repository server is an internal service which maintains a local cache of the Git repository holding the application manifests.  
It is responsible for generating and returning rendered Kubernetes manifests when provided the following input:

1. Repository URL
2. Revision (commit, tag, branch)
3. Application path
4. Template‑specific settings: parameters, Helm `values.yaml`, etc.

---

#### Application Controller ⚙️

The application controller is a Kubernetes controller which continuously monitors running applications and compares the current
live state against the desired target state (as specified in the repo).  
It detects **Out‑of‑Sync** state and (optionally) takes corrective action, and is responsible for invoking any user‑defined hooks for lifecycle events (`PreSync`, `Sync`, `PostSync`).

---

#### ApplicationSet Controller ♾️

The **ApplicationSet Controller** (`argocd-applicationset-controller`) watches `ApplicationSet` CRDs and dynamically **generates or prunes `Application` resources** based on template + generator patterns (cluster list, Git directories, pull requests, matrices, etc.).  
This enables *fleet‑style* GitOps where dozens or hundreds of Applications are managed from a single declarative spec.

---

#### Redis ⚡

`argocd-redis` (or `argocd-redis-ha`) provides a shared, in‑memory data store used for:

* Caching application trees and diff results
* Persisting UI session tokens
* Streaming events between controllers (HA mode)

In single‑node installs Argo CD can fall back to on‑disk BoltDB, but HA deployments **require** Redis for consistency and performance.

---

#### Dex Server 🔐

`argocd-dex-server` bundles **Dex**, an OpenID Connect identity provider which brokers authentication to external IdPs (GitHub, GitLab, LDAP, SAML, Azure AD, …).  
If you already have an enterprise IdP you can disable Dex and point Argo CD directly at your own OIDC endpoint.

---

#### Notifications Controller 📣

`argocd-notifications-controller` watches application events (sync‑succeeded, health‑degraded, etc.) and delivers messages via:

* Slack / Microsoft Teams
* Email / SMTP
* Generic Webhooks / SNS / Opsgenie / PagerDuty

Notification rules, triggers, and message templates are defined in two ConfigMaps:  
`argocd-notification-templates-cm` and `argocd-notification-trigger-cm`.

---

## 🔧 支持的 Git 类型与模板工具

### 支持的 Git 仓库类型：

- ✅ GitHub / GitHub Enterprise
- ✅ GitLab / GitLab Self-Hosted
- ✅ Bitbucket / Bitbucket Server
- ✅ Gitea、Azure DevOps 等


## 📚 推荐阅读

- GitOps 官方定义：[https://opengitops.dev](https://opengitops.dev)
- Argo CD 架构文档：[https://argo-cd.readthedocs.io/en/stable/operator-manual/architecture/](https://argo-cd.readthedocs.io/en/stable/operator-manual/architecture/)
- GitOps vs CI/CD 博客：https://www.weave.works/blog/gitops-operations-by-pull-request

