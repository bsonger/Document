---
title: "如何创建kubeConfig"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# 如何创建kubeConfig

1. 创建clusterRole

> kubectl create clusterrole monitor --resource=prometheusrule --verb="*"

2. 创建clusterRoleBinding
> kubectl create clusterrolebinding monitor --clusterrole=rule --user monitor
3. 生成普通用户令牌
```shell
cd /etc/kubernetes/pki #你的集群证书目录
 
 
user=monitor
organization=monitoring
clustername=cluster02
# 【集群地址，例如：https://172.21.114.169:6443】
api_addr=https://192.168.31.208:6443
time_days=365
 
#复制粘贴下面的代码即可
 
umask 077;openssl genrsa -out $user.key 2048
 
openssl req -new -key $user.key -out $user.csr -subj "/O=$organization/CN=$user"
 
openssl  x509 -req -in $user.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $user.crt -days $time_days
 
kubectl config set-cluster $clustername --server=$api_addr --certificate-authority=ca.crt --embed-certs=true --kubeconfig=/root/$user.config
 
kubectl config set-credentials $user --client-certificate=$user.crt --client-key=$user.key --embed-certs=true --kubeconfig=/root/$user.config
 
kubectl config set-context $user@$clustername --cluster=$clustername --user=$user --kubeconfig=/root/$user.config
 
kubectl config use-context $user@$clustername --kubeconfig=/root/$user.config
```