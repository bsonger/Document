---
title: "创建一个无状态服务"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---
## 创建一个deployment资源

**当创建一个deployment时，kubernetes 内部是如何处理的？**

> 1. 创建 Deployment 资源
>   - 用户通过 kubectl apply 或 API 直接创建一个 Deployment 资源。
>   - kube-apiserver 接收到请求后，将 Deployment 对象写入 etcd 数据库。
>2. controller-manager 处理 Deployment
>   - kube-controller-manager 监听 Deployment 资源（通过 watch 机制）。
>   - 发现新建的 Deployment，DeploymentController 计算出需要创建的 ReplicaSet（如果是更新，还会进行滚动升级）。
>   - controller-manager 通过 kube-apiserver 创建 ReplicaSet 资源。
>3. controller-manager 处理 ReplicaSet
>   - kube-controller-manager 监听 ReplicaSet 资源（通过 watch 机制）。
>   - 发现新建的 ReplicaSet，ReplicaSetController 计算出需要创建的 Pod。
>   - controller-manager 通过 kube-apiserver 创建 Pod 资源。
>4. kube-scheduler 进行调度
>   - kube-scheduler 监听 Pod 资源的变更（通过 watch 机制）。
>   - 发现 Pod 处于 Pending 状态，调度器计算出合适的 Node（根据 CPU、内存、调度策略等）。
>   - kube-scheduler 通过 kube-apiserver 更新 Pod 的 spec.nodeName，分配到合适的节点。
>5. kubelet 负责拉起容器
>   - kubelet 监听 kube-apiserver，发现被分配到自己节点的 Pod。
>   - kubelet 通过 Container Runtime Interface (CRI) 执行：
>     - 创建网络（调用 CNI 插件，分配 IP）。
>     - 挂载存储（调用 CSI 插件，挂载 Volume）。
>     - 启动容器（调用 CRI，使用 containerd 或 Docker 等创建容器）。
>   - kubelet 通过 kube-apiserver 更新 Pod 状态（Running）。

## 创建一个service资源
> 1. kube-proxy 监听 kube-apiserver 上的 service 资源的创建：
>   - 当 kube-proxy 监视到新的服务资源（Service）时，它会根据 Service 中的标签（labels）来筛选出符合条件的 Pod。这些 Pod 是与该 Service 相关联的目标（后端）。通过这些信息，kube-proxy 创建一个 Endpoints 资源对象，指向这些 Pod。
> 2. kube-proxy 监听 endpoint 资源的变化：
>   - Endpoints 资源对象被创建或更新，kube-proxy 就会动态地根据新的 Endpoint 信息来修改其网络规则（如 iptables 或 ipvs）。
>   - kube-proxy 根据 Endpoint 中的 Pod 地址创建或更新相应的 iptables 或 ipvs 规则，这样流量就可以通过正确的后端 Pod 进行路由。