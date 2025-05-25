---
title: "堆"
weight: 5
bookCollapseSection: true
---

本章节介绍堆

## 基本思想
- 堆通常用完全二叉树实现，并存储在数组中，保证高效的插入和删除操作
- 最大堆：每个节点的值都大于等于其子节点的值，根节点是最大值
- 最小堆：每个节点的值都小于其子节点的值，根节点是最小值

## 存储方式
- 使用数组表示完全二叉树

## 操作
插入操作
1. 把新元素插入到数组末尾
2. 上浮

删除堆顶元素
1. 用最后一个元素替换堆顶
2. 下沉

堆排序
1. 构建最大堆
2. 交换根节点和最后一个元素
3. 调整堆
4. 重复上述步骤，直到排序完成


## 经典应用
- 优先队列（Priority Queue）
  - 使用最小堆实现任务优先级调度
  - 使用最大堆实现Top K 问题
- 求TOP K 大/小元素
  - 最小堆求TOP K 大元素
  - 最大堆求TOP K 小元素
- 中位数维护
- 左边最大堆（维护较小的一半数）
- 右边最小堆（维护较大的一半数）
- 中位数为两堆的顶元素

## 优缺点
优点
- 插入、删除 操作时间复杂度 O(log n)
- 动态维护 Top K 问题，适用于大数据场景
- 优先队列实现简单高效

缺点
- 不支持快速查找（不像哈希表 O(1)）
- 堆排序不是稳定排序（相同元素相对顺序可能改变）
- 数组实现堆可能需要额外的空间移动

## 示例

最小堆
```go
import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] } // 最小堆
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func main() {
    h := &MinHeap{}
    heap.Init(h)
    heap.Push(h, 5)
    heap.Push(h, 2)
    heap.Push(h, 8)
    heap.Push(h, 1)
    fmt.Println(heap.Pop(h)) // 1
}
```