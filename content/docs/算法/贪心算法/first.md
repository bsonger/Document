---
title: "跳跃游戏 2"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# 跳跃游戏 2 
[leetcode 45](https://leetcode.cn/problems/jump-game-ii/?envType=study-plan-v2&envId=top-100-liked "跳跃游戏 2")

给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。

每个元素 nums[i] 表示从索引 i 向后跳转的最大长度。换句话说，如果你在 nums[i] 处，你可以跳转到任意 nums[i + j] 处:

0 <= j <= nums[i]
i + j < n
返回到达 nums[n - 1] 的最小跳跃次数。生成的测试用例可以到达 nums[n - 1]。



示例 1:

输入: nums = [2,3,1,1,4]
输出: 2
解释: 跳到最后一个位置的最小跳跃数是 2。
从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。
示例 2:

输入: nums = [2,3,0,1,4]
输出: 2


```go

func jump(nums []int) int {
    n := len(nums)
    if n <= 1{
        return 0
    }
    steps , curEnd, maxReach := 0, 0, 0
    for i :=0 ; i< n ; i ++{
        maxReach = max(maxReach, i + nums[i])
        if i == curEnd{
            steps ++
            curEnd = maxReach
            if curEnd >= n-1{
                break
            }
        }
    }
    return steps
}

```
