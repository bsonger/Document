---
title: "拓扑排序"
weight: 6
bookCollapseSection: true
---

本章节介绍堆

## 基本思想
拓扑排序（Topological Sorting）是 针对有向无环图（DAG, Directed Acyclic Graph）的一种排序方法，它的目标是将图中的所有顶点排序，使得对于每一条有向边 (u → v)，顶点 u 先出现在 v 之前。

1. kahn算法（入度法）
   - 维护所有入度0的节点（可以入队）
   - 每次从队列中取出一个节点，并删除它指向的边，更新相邻节点的入度
   - 若某个相邻节点入度变成0，则加入队列
   - 知道所有节点都被访问完
2. DFS 逆序法
   - 采用DFS遍历每个点，递归访问所有子节点
   - 访问完所有后续节点后，将当前的节点压入栈
   - 最终栈中的元素顺序即为拓扑排序的结果

## 经典应用
1. 任务调度
- 场景：多个任务有前后依赖关系，必须按顺序执行。
- 示例：编译多个源代码文件时，需要先编译依赖文件。

2. 课程安排（先修课程）
- 场景：某些课程需要先完成前置课程，才能学习后续课程。
- 示例：LeetCode 207 课程表（Course Schedule）问题。

3. 电路时序分析
- 场景：计算电路中信号的传播顺序，确保数据按逻辑流动。

4. 软件包依赖管理
- 场景：如 npm、pip 等包管理器，安装包时需要先安装其依赖项。

5. 编译器优化
- 场景：代码中的函数调用优化，确定依赖关系，减少重复计算。
## 优缺点
优点

- 高效性：
  - Kahn 算法 时间复杂度 O(V + E)（遍历所有节点和边）。
  - DFS 逆序法 时间复杂度也是 O(V + E)。
- 易于实现：
  - 代码简单，适用于有向无环图（DAG）。
- 解决依赖关系问题：
  - 能很好地解决 任务调度、先修课程、软件包管理 等依赖问题。

缺点

- 仅适用于 DAG：
  - 不能处理有环图，如果图中有环，无法进行拓扑排序。
- 可能有多种拓扑排序结果：
  - 拓扑排序的结果可能不是唯一的，不同的算法可能得出不同的合法顺序。
- 不能检测所有依赖冲突：
  - 仅能发现 循环依赖（环），但无法处理复杂的冲突（如并发依赖）。