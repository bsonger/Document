---
title: "Kubernetes"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: true
---
kubernetes 架构图


![kubernetes 架构图](/images/kubernetes.png)


## Group、Version、Resource 核心数据结构

整个Kubernetes 体系架构中，**资源**是Kubernetes最重要的概念，可以说Kubernets的生态系统都围绕着资源运作。

Kubernetes系统虽然有相当复杂和众多的功能，但是它本质上是一个资源控制系统——注册、管理、调度资源并维护资源的状态

Kubernetes 将资源分组和版本化，形成Group（资源组）、Version（资源版本）、Resource（资源）

- **Group：资源组，在Kubernetes API Server 中也可称其为APIGroup**
- **Version：资源版本，在Kubernets API Server 中也可称其为APIVersions**
- **Resource：资源，在Kubernetes API Server 中也可称其为APIResource**
- **Kind：资源种类，描述Resource 的种类，与Resource为同一级别**

Kubernetes 系统支持多个Group。每个Group支持多个Version，每个Version支持多个Resource，同时部分资源会拥有自己的子资源

资源组、资源版本、资源、子资源的完整表达形式为\<group>/\<version>/\<resource>/\<status>,以常用的Deployment资源为例，其完整表现形式为apps/v1/deployments/status

资源对象由 资源组 + 资源版本 + 资源种类组成，例如Deployment 资源是实例化后拥有资源组、资源版本及资源种类，其表现形式为<group>/<version>,Kind=<kind>,例如apps/v1，Kind= Deployment

每一个资源都拥有一定数量的资源操作方法（即Verbs），资源操作方法用于Etcd集群存储中对资源对象的增、删、改、查操作。目前Kubernetes系统支持8种资源操作方法，分别是create、delete、deletecollection、get、list、patch、update、watch操作方法

