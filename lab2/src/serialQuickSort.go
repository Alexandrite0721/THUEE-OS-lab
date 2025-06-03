package main

import "math/rand"

// 串行快速排序的实现
func serialQuickSort(arr []int, low, high int) {
	if low < high {
		// 获取分区点
		pivotIndex := partition(arr, low, high)

		// 递归排序左右两部分
		serialQuickSort(arr, low, pivotIndex-1)
		serialQuickSort(arr, pivotIndex+1, high)
	}
}

// 快速排序的分区函数
func partition(arr []int, low, high int) int {
	// 随机选择一个元素作为pivot，避免最坏情况
	pivotIdx := low + rand.Intn(high-low+1)
	arr[pivotIdx], arr[high] = arr[high], arr[pivotIdx]

	// 选择最后一个元素为pivot
	pivot := arr[high]

	// i 是小于pivot的元素应该放的位置
	i := low - 1

	for j := low; j < high; j++ {
		// 如果当前元素小于pivot
		if arr[j] < pivot {
			i++
			// 交换arr[i]和arr[j]
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// 将pivot放在正确的位置
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
