package main

import (
	"fmt"
	"strings"
)

func main() {
	// 测试用例1: 安全分配
	fmt.Println("===== 测试用例 1: 安全分配 =====")
	system1 := NewSystem(
		5,              // 进程数
		3,              // 资源种类
		[]int{3, 3, 2}, // 可用资源
		// 最大需求矩阵
		[][]int{
			{7, 5, 3},
			{3, 2, 2},
			{9, 0, 2},
			{2, 2, 2},
			{4, 3, 3},
		},
		// 已分配矩阵
		[][]int{
			{0, 1, 0},
			{2, 0, 0},
			{3, 0, 2},
			{2, 1, 1},
			{0, 0, 2},
		},
	)
	system1.PrintState()
	system1.IsSafe()

	// P1请求资源 (1,0,2)
	request1 := []int{1, 0, 2}
	system1.RequestResource(1, request1)
	system1.PrintState()

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试用例2: 不安全请求
	fmt.Println("\n===== 测试用例 2: 不安全分配 =====")
	system2 := NewSystem(
		5,
		3,
		[]int{3, 3, 2},
		[][]int{
			{7, 5, 3},
			{3, 2, 2},
			{9, 0, 2},
			{2, 2, 2},
			{4, 3, 3},
		},
		[][]int{
			{0, 1, 0},
			{2, 0, 0},
			{3, 0, 2},
			{2, 1, 1},
			{0, 0, 2},
		},
	)
	system2.PrintState()

	// P0请求资源 (0,3,0) - 会导致不安全状态
	request2 := []int{0, 3, 0}
	system2.RequestResource(0, request2)
	system2.PrintState()

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试用例3: 超过需求
	fmt.Println("\n===== 测试用例 3: 超过最大需求 =====")
	system3 := NewSystem(
		5,
		3,
		[]int{3, 3, 2},
		[][]int{
			{7, 5, 3},
			{3, 2, 2},
			{9, 0, 2},
			{2, 2, 2},
			{4, 3, 3},
		},
		[][]int{
			{0, 1, 0},
			{2, 0, 0},
			{3, 0, 2},
			{2, 1, 1},
			{0, 0, 2},
		},
	)
	system3.PrintState()

	// P1请求资源 (2,1,1) - 超过最大需求
	request3 := []int{2, 1, 1}
	system3.RequestResource(1, request3)
	system3.PrintState()
}
