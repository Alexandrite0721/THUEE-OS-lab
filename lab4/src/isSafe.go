package main

import "fmt"

// IsSafe 安全性检查
func (s *System) IsSafe() bool {
	fmt.Println("\n开始安全性检查...")

	// 初始化工作向量和完成标记
	work := make([]int, s.resources)
	finish := make([]bool, s.processes)
	safeSequence := make([]int, 0, s.processes)

	copy(work, s.available)

	fmt.Printf("初始工作向量: %v\n", work)

	found := true
	for found {
		found = false

		for i := 0; i < s.processes; i++ {
			if !finish[i] {
				// 检查进程i的需求是否小于等于工作向量
				canRun := true
				for j := 0; j < s.resources; j++ {
					if s.need[i][j] > work[j] {
						canRun = false
						break
					}
				}

				if canRun {
					fmt.Printf("  - 进程 P%d 可执行 (需求: %v <= 工作向量: %v)\n", i, s.need[i], work)

					// 模拟进程执行完成
					for j := 0; j < s.resources; j++ {
						work[j] += s.allocation[i][j]
					}

					finish[i] = true
					safeSequence = append(safeSequence, i)
					found = true

					fmt.Printf("    进程 P%d 完成, 释放资源: %v\n", i, s.allocation[i])
					fmt.Printf("    更新工作向量: %v\n", work)
				}
			}
		}
	}

	// 检查是否所有进程都已完成
	allFinished := true
	for i := 0; i < s.processes; i++ {
		if !finish[i] {
			allFinished = false
			break
		}
	}

	if allFinished {
		fmt.Printf("安全序列: ")
		for i, pid := range safeSequence {
			fmt.Printf("P%d", pid)
			if i < len(safeSequence)-1 {
				fmt.Print(" -> ")
			}
		}
		fmt.Println("\n系统处于安全状态")
		return true
	} else {
		fmt.Println("未找到安全序列! 系统处于不安全状态")
		return false
	}
}
