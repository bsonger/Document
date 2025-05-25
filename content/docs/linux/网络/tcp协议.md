---
title: "tcp 协议"
weight: 1
bookCollapseSection: false
---

## 详解三次握手，四次挥手

#### 三次握手（建立连接）

TCP采用 三次握手（Three-Way Handshake） 过程来建立可靠连接：

1. 客户端发送 SYN（同步序列号）：
   - 客户端向服务器发送一个 SYN 报文，表示请求建立连接。
   - 该报文的 SYN=1，并包含一个初始序列号 seq=x。 
2. 服务器回应 SYN-ACK：
   - 服务器收到 SYN 后，发送 SYN=1, ACK=1，并附带自己的初始序列号 seq=y，确认客户端的 seq=x+1。
3. 客户端回应 ACK：
   - 客户端收到服务器的 SYN-ACK 后，发送 ACK=1，确认 seq=y+1，同时携带 seq=x+1。

#### 四次挥手（断开连接）

TCP 采用 四次挥手（Four-Way Handshake） 断开连接：
1. 客户端发送 FIN：
   - 客户端想要断开连接，发送 FIN=1，表示无数据可发。
2. 服务器发送 ACK：
   - 服务器收到 FIN 后，回复 ACK=1，但可能仍有未发送的数据。
3. 服务器发送 FIN：
   - 服务器确认数据发送完毕，发送 FIN=1 断开连接。
4. 客户端发送 ACK，完成断开：
   - 客户端确认 FIN，回复 ACK=1，连接正式关闭。
