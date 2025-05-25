---
title: "configuration"
weight: 3
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---


A single-page markdown reference that condenses:

1. **Global / atomic config objects**
2. **Multi‑object layout (Application + Project + repo Secret)**
3. **Application 🛰️ essentials** – full **`syncPolicy`** deep‑dive
4. **AppProject 🏗️ guard‑rails**
5. **ApplicationSet ♾️ mass‑generation**

---

## 1  🗄️ Core (“Atomic”) Configuration Objects
| Sample file | K8s name | Kind | Purpose |
|---|---|---|---|
| `argocd-cm.yaml` | argocd‑cm | ConfigMap | Global Argo CD behaviour (logo, OIDC, diff‑settings). |
| `argocd-repositories` | _*_ | Secret | Git / Helm repo endpoints + creds. |
| `argocd-repo-creds.yaml` | _*_ | Secret | Credential **templates** shared by many repos. |
| `argocd-cmd-params-cm.yaml` | argocd-cmd-params-cm | ConfigMap | Extra env / CLI flags for all components. |
| `argo-cd-secret.yaml` | argocd-secret | Secret | Admin pwd, JWT key, Dex & webhook secrets. |
| `argocd-rbac-cm.yaml` | argocd-rbac-cm | ConfigMap | Casbin RBAC rules for UI/CLI. |
| `argocd-tls-certs-cm.yaml` | argocd-tls-certs-cm | ConfigMap | Custom CA bundles for outbound HTTPS. |
| `argocd-ssh-known-hosts-cm.yaml` | argocd-ssh-known-hosts-cm | ConfigMap | `known_hosts` for outbound SSH. |

---

## 2  🗂️ Multi‑Object Layout

| File | Kind | Description |
|---|---|---|
| `application.yaml` | Application 🛰️ | One workload stack. |
| `project.yaml` | AppProject 🏗️ | Guard‑rails for that workload(s). |
| `argocd-repositories` | Secret 🔐 | Git token / SSH key for the repo. |

---

## 3  🛰️ Application – Key Fields & `syncPolicy`

### Minimal skeleton
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: sample
spec:
  project: default
  source:
    repoURL: https://github.com/org/repo.git
    path: k8s/app
    targetRevision: main
  destination:
    server: https://kubernetes.default.svc
    namespace: sample
```

### `syncPolicy` cheat‑sheet
```yaml
syncPolicy:
  automated:        # ⚡ auto‑deploy
    prune: true
    selfHeal: true
  syncOptions:      # 🛡️ how kubectl behaves
    - CreateNamespace=true
    - ApplyOutOfSyncOnly=true
    - PruneLast=true
  retry:            # 🔁 back‑off
    limit: 3
    backoff: {duration: 10s, factor: 2, maxDuration: 2m}
  managedNamespaceMetadata:  # 🏷️ labels on auto‑ns
    labels: {owner: platform}
```

#### Full `syncOptions` catalogue

| Option | What it does |
|---|---|
| `CreateNamespace=true` | Auto‑create `destination.namespace`. |
| `ApplyOutOfSyncOnly=true` | Patch only drifted resources. |
| `PruneLast=true` | Run prune step after all others. |
| `PrunePropagationPolicy=<mode>` | Foreground / background / orphan delete. |
| `Replace=true` | Use *replace* instead of *apply*. |
| `Force=true` | Delete & recreate conflicting objects. |
| `ServerSideApply=true` | Use SSA; avoids 256 KiB anno limit. |
| `Validate=false` | Skip schema validation. |
| `SkipDryRunOnMissingResource=true` | Skip dry‑run for CRDs not yet installed. |
| `RespectIgnoreDifferences=true` | Honour ignore rules during sync. |
| `FailOnSharedResource=true` | Abort if resource managed by another App. |
| `allowEmpty=true` *(automated)* | Sync even if path empty (v2.4+). |

---

## 4  🏗️ AppProject – Guard‑Rails

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: team-a
spec:
  description: Team‑A sandbox
  sourceRepos: [ 'https://github.com/org/*' ]   # allowed repos
  destinations:
    - {server: https://kubernetes.default.svc, namespace: team-a}
  clusterResourceWhitelist: [{group: '', kind: Namespace}]
  namespaceResourceBlacklist:
    - {group: '', kind: LimitRange}
    - {group: '', kind: NetworkPolicy}
  roles:
    - name: dev
      groups: [ team-a-devs ]                   # SSO group
      policies:
        - p, proj:team-a:dev, applications, sync, team-a/*, allow
```

---

## 5  ♾️ ApplicationSet – Mass Generation

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: guestbook
spec:
  generators:
    - list:
        elements:
          - {cluster: dev,  url: https://1.2.3.4}
          - {cluster: prod, url: https://2.4.6.8}
  template:
    metadata:
      name: '{{.cluster}}-guestbook'
    spec:
      project: team-a
      source:
        repoURL: https://github.com/org/clusters.git
        path: guestbook/{{.cluster}}
        targetRevision: HEAD
      destination:
        server: '{{.url}}'
        namespace: guestbook
```

---

## 6  🧠 TL;DR Matrix

| Resource | Scope | Core Job |
|---|---|---|
| Application 🛰️ | 1 workload | Sync Git → cluster/namespace. |
| AppProject 🏗️ | many workloads | Define repos, target clusters, RBAC. |
| ApplicationSet ♾️ | many Applications | Auto‑stamp Apps from patterns. |
| \*‑cm/secret | cluster‑global | Tune Argo CD engine & credentials. |

---

**Version note:** verified against Argo CD **v2.6 +**.  
Copy → save → GitOps on! ✨
