package main

import "fmt"

// RequestResource 请求资源
func (s *System) RequestResource(pid int, request []int) bool {
	fmt.Printf("\n===== 进程 P%d 请求资源: %v =====\n", pid, request)

	// 步骤1: 检查请求是否超过需求
	for j := 0; j < s.resources; j++ {
		if request[j] > s.need[pid][j] {
			fmt.Printf("错误: 请求超过需求 (资源%d: 请求%d > 需求%d)\n", j, request[j], s.need[pid][j])
			return false
		}
	}

	// 步骤2: 检查请求是否超过可用资源
	for j := 0; j < s.resources; j++ {
		if request[j] > s.available[j] {
			fmt.Printf("错误: 请求超过可用资源 (资源%d: 请求%d > 可用%d)\n", j, request[j], s.available[j])
			return false
		}
	}

	// 步骤3: 尝试分配资源
	fmt.Println("尝试分配资源...")
	backupAvailable := make([]int, s.resources)
	backupAllocation := make([]int, s.resources)
	backupNeed := make([]int, s.resources)

	copy(backupAvailable, s.available)
	copy(backupAllocation, s.allocation[pid])
	copy(backupNeed, s.need[pid])

	for j := 0; j < s.resources; j++ {
		s.available[j] -= request[j]
		s.allocation[pid][j] += request[j]
		s.need[pid][j] -= request[j]
	}

	// 步骤4: 安全性检查
	if s.IsSafe() {
		fmt.Printf("资源分配成功! 系统处于安全状态\n")
		return true
	} else {
		// 恢复原始状态
		fmt.Println("资源分配导致不安全状态，撤销分配")
		copy(s.available, backupAvailable)
		copy(s.allocation[pid], backupAllocation)
		copy(s.need[pid], backupNeed)
		return false
	}
}
