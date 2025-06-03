package main

import "sync"

// 并行快速排序函数
func parallelQuickSort(arr []int, low, high int, wg *sync.WaitGroup) {
	// 如果数据量小，使用串行快速排序
	if high-low < minSize {
		serialQuickSort(arr, low, high)
		return
	}

	// 如果数据范围有效，则进行排序
	if low < high {
		// 分区，获取分区点
		pivot := partition(arr, low, high)

		// 检查是否有可用的工作线程来并行处理左半部分
		var canParallel bool

		// 检查是否可以创建新的工作线程
		workerMutex.Lock()
		if activeWorkers < maxWorkers {
			activeWorkers++
			canParallel = true
		}
		workerMutex.Unlock()

		if canParallel {
			// 并行处理左半部分
			wg.Add(1)
			go func() {
				defer wg.Done()
				parallelQuickSort(arr, low, pivot-1, wg)

				// 工作完成后减少活动工作线程计数
				workerMutex.Lock()
				activeWorkers--
				workerMutex.Unlock()
			}()

			// 当前线程处理右半部分
			parallelQuickSort(arr, pivot+1, high, wg)
		} else {
			// 在当前线程中顺序处理两部分
			parallelQuickSort(arr, low, pivot-1, wg)
			parallelQuickSort(arr, pivot+1, high, wg)
		}
	}
}
