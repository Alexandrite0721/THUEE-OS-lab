package main

import "fmt"

// System 系统资源状态
type System struct {
	processes  int     // 进程数量
	resources  int     // 资源种类数量
	available  []int   // 可用资源向量
	max        [][]int // 最大需求矩阵
	allocation [][]int // 分配矩阵
	need       [][]int // 需求矩阵
}

// NewSystem 初始化系统
func NewSystem(processes, resources int, available []int, max, allocation [][]int) *System {
	system := &System{
		processes:  processes,
		resources:  resources,
		available:  make([]int, resources),
		max:        make([][]int, processes),
		allocation: make([][]int, processes),
		need:       make([][]int, processes),
	}

	copy(system.available, available)

	for i := 0; i < processes; i++ {
		system.max[i] = make([]int, resources)
		system.allocation[i] = make([]int, resources)
		system.need[i] = make([]int, resources)
		copy(system.max[i], max[i])
		copy(system.allocation[i], allocation[i])

		// 计算需求矩阵
		for j := 0; j < resources; j++ {
			system.need[i][j] = max[i][j] - allocation[i][j]
		}
	}

	return system
}

// PrintState 打印系统状态
func (s *System) PrintState() {
	fmt.Println("\n当前系统状态:")
	fmt.Println("可用资源:", s.available)

	fmt.Println("\n进程\t最大需求\t已分配\t需求")
	for i := 0; i < s.processes; i++ {
		fmt.Printf("P%d\t%v\t%v\t%v\n", i, s.max[i], s.allocation[i], s.need[i])
	}
}
