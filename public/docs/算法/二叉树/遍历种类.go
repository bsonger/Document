package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func preOrderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	stack, res := []*TreeNode{root}, []int{}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, node.Val)
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
	return res
}

func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	stack, res := []*TreeNode{}, []int{}
	cur := root
	for cur != nil && len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		cur = root.Right
	}
	return res
}

func subsets(nums []int) [][]int {
	var result [][]int
	var path []int
	var dfs func(start int)
	dfs = func(start int) {
		tmp := make([]int, len(path))
		copy(tmp, path)
		result = append(result, tmp)
		for i := start; i < len(nums); i++ {
			path = append(path, nums[i])
			dfs(i + 1)
			path = path[:len(path)-1]
		}
	}
	dfs(0)
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(subsets(nums))
}
