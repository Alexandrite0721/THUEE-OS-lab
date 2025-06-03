package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
)

// 生成随机数文件
func generateRandomFile() error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < totalNumbers; i++ {
		num := rand.Intn(totalNumbers * 10) // 生成0~10M的随机数
		_, err := writer.WriteString(strconv.Itoa(num) + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// 从文件读取数据
func readDataFromFile() ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		data = append(data, num)
	}
	return data, scanner.Err()
}

// 写入排序结果到文件
func writeDataToFile(data []int) error {
	file, err := os.Create(sortedFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, num := range data {
		_, err := writer.WriteString(strconv.Itoa(num) + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
