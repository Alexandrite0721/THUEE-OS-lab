package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	numberMutex   sync.Mutex
	currentNumber = 0
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <num_clerks>")
		return
	}
	n, _ := strconv.Atoi(os.Args[1])

	scanner := bufio.NewScanner(os.Stdin)
	var customers []Customer
	fmt.Println("请输入顾客信息，每行格式为：顾客序号 进入银行时间 服务时间，输入空行结束：")
	for {
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			break
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Printf("该行输入格式不正确：%s，请重新输入\n", line)
			continue
		}
		id, _ := strconv.Atoi(fields[0])
		arrival, _ := strconv.Atoi(fields[1])
		service, _ := strconv.Atoi(fields[2])
		customers = append(customers, Customer{
			CustomerID:  id,
			ArrivalTime: arrival,
			ServiceTime: service,
		})
	}

	sort.Slice(customers, func(i, j int) bool {
		return customers[i].ArrivalTime < customers[j].ArrivalTime
	})

	q := &CustomerQueue{
		customers: make([]*Customer, 0),
	}
	q.cond = sync.NewCond(&q.mutex)

	results := make(chan *Customer, 100)

	// 启动柜员协程
	clerkWg := startClerkWorkers(n, q, results)
	// 启动顾客协程
	customerWg := startCustomerWorkers(customers, q)

	customerWg.Wait()

	q.mutex.Lock()
	for i := 0; i < n; i++ {
		q.customers = append(q.customers, nil)
	}
	q.cond.Broadcast()
	q.mutex.Unlock()

	clerkWg.Wait()
	close(results)

	var resultList []*Customer
	for c := range results {
		if c != nil {
			resultList = append(resultList, c)
		}
	}

	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].CustomerID < resultList[j].CustomerID
	})

	for _, c := range resultList {
		fmt.Printf("顾客id：%d 到达时间：%d 开始服务时间：%d 结束服务时间：%d 服务柜员id：%d\n",
			c.CustomerID, c.ArrivalTime, c.StartTime, c.EndTime, c.ClerkID)
	}
}
