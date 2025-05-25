---
title: "kubernetes service"
weight: 1
#bookCollapseSection: true
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      # 提取 annotation 中的 metrics_path
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__

      # 提取 annotation 中的 scrape interval
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape_interval]
        action: replace
        target_label: scrape_interval

      # 提取 annotation 中的 scrape timeout
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape_timeout]
        action: replace
        target_label: scrape_timeout

      # 通过 annotation 设置是否启用 scrape
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: "true"

      # 提取 annotation 中的 metrics 端口
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: (.+)
        replacement: $1
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-app
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/path: "/metrics"
    prometheus.io/port: "8080"
    prometheus.io/scrape_interval: "15s"
    prometheus.io/scrape_timeout: "10s"
spec:
  containers:
    - name: example
      image: example-app
      ports:
        - containerPort: 8080
```