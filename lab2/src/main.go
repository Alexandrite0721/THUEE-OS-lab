package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	totalNumbers = 1000000
	minSize      = 1000
	maxWorkers   = 20
	filename     = "random_numbers.txt"
	sortedFile   = "sorted_numbers.txt"
)

var (
	activeWorkers int
	workerMutex   sync.Mutex
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 生成随机数文件
	if err := generateRandomFile(); err != nil {
		fmt.Println("生成随机数文件时出错:", err)
		return
	}
	fmt.Println("随机数文件已准备")

	// 读取数据文件
	data, err := readDataFromFile()
	if err != nil {
		fmt.Println("读取数据文件时出错:", err)
		return
	}
	fmt.Printf("已读取 %d 个数字\n", len(data))

	// 开始排序计时
	startTime := time.Now()
	fmt.Println("开始并行排序...")

	// 创建排序数据的副本并排序
	sortedData := make([]int, len(data))
	copy(sortedData, data)

	// 创建一个等待组，用于等待所有排序任务完成
	var wg sync.WaitGroup
	wg.Add(1)

	// 开始并行快速排序
	go func() {
		defer wg.Done()
		// 获取工作线程锁并增加活动工作线程计数
		workerMutex.Lock()
		activeWorkers++
		workerMutex.Unlock()

		// 启动递归并行排序
		parallelQuickSort(sortedData, 0, len(sortedData)-1, &wg)

		// 完成后减少活动工作线程计数
		workerMutex.Lock()
		activeWorkers--
		workerMutex.Unlock()
	}()

	// 等待所有排序任务完成
	wg.Wait()

	// 结束计时并打印信息
	duration := time.Since(startTime)
	fmt.Printf("排序完成，耗时: %v\n", duration)

	// 将排序结果写入文件
	if err := writeDataToFile(sortedData); err != nil {
		fmt.Println("写入排序结果时出错:", err)
		return
	}
	fmt.Printf("排序结果已保存至文件: %s\n", sortedFile)

	// 验证排序结果
	if isSorted(sortedData) {
		fmt.Println("排序结果验证: 成功！数据已正确排序")
	} else {
		fmt.Println("排序结果验证: 失败！数据未正确排序")
	}
}
