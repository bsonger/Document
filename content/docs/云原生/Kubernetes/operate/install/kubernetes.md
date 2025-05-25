---
title: "kubernetes"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
#bookCollapseSection: true
---

## 部署前准备




1. 修改机器名
```shell
hostnamectl set-hostname master
````
2. 关闭swap
```shell
swapoff -a
sed -i '/swap / s/^\(.*\)$/#\1/g' /etc/fstab
free -m # 检查是否关闭swap

cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

sudo modprobe br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system
````
3. ssh允许Root 登陆
```shell
vim /etc/ssh/sshd_config

增加

Port 22
PermitRootLogin yes

passwd # 设置root 密码
````
4. 关闭防火墙

```shell
sudo ufw disable
````

5. 修改apt 源地址

```shell
cp /etc/apt/sources.list /etc/apt/sources.list.bak
# 使用vim 替换/etc/apt/sources.list中的资源地址

vim /etc/apt/sources.list

deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
# 更新
vim /etc/resolv.conf
# 添加一下两行
nameserver 8.8.8.8
nameserver 8.8.4.4

apt-get update

````
6.  添加 Kubernetes apt 存储库
```shell
sudo tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main
EOF

apt-get update



The following signatures couldn't be verified because the public key is not available: NO_PUBKEY B53DC80D13EDEF05 NO_PUBKEY FEEA9169307EA071
# 如果报错以上错误，执行下面一条命令 recv-keys 是报错现实的NO_PUBKEY 
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys FEEA9169307EA071
```
7. 安装kubectl kubeadm kubelet
```shell
sudo apt-get install -y kubelet=1.22.2-00 kubeadm=1.22.2-00 kubectl=1.22.2-00
```
8. 安装docker
```shell
apt-get install -y docker.io
systemctl enable docker
```

## master

1.  安装master 节点

```shell
kubeadm init \
--apiserver-advertise-address=192.168.31.239 \
--image-repository registry.aliyuncs.com/google_containers \
--kubernetes-version v1.32.3 \
--service-cidr=10.96.0.0/12 \
--pod-network-cidr=10.244.0.0/16
```

## node
1. 安装node 节点
```shell
kubeadm join 192.168.31.208:6443 --token u1soc8.1xyeqbzptpvvjz8f     --discovery-token-ca-cert-hash sha256:3406b992b2c05f27b398a81375082a72aa2a823fd82616e4f8cb6f2a24370bb3
# kubeadm init 初始化完成会打印到console 上
# 可以通过 kubeadm token create --print-join-command 查看
```