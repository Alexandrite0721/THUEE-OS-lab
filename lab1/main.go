package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Customer struct {
	CustomerID  int
	ArrivalTime int
	ServiceTime int
	StartTime   int
	EndTime     int
	ClerkID     int
}

type Queue struct {
	customers []*Customer
	mutex     sync.Mutex
	cond      *sync.Cond
}

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

	q := &Queue{
		customers: make([]*Customer, 0),
	}
	q.cond = sync.NewCond(&q.mutex)

	results := make(chan *Customer, 100)

	var clerkWg sync.WaitGroup
	clerkWg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			defer clerkWg.Done()
			clerkAvailableTime := 0
			for {
				q.mutex.Lock()
				for len(q.customers) == 0 {
					q.cond.Wait()
				}
				customer := q.customers[0]
				q.customers = q.customers[1:]
				q.mutex.Unlock()

				if customer == nil {
					return
				}

				startTime := max(customer.ArrivalTime, clerkAvailableTime)
				endTime := startTime + customer.ServiceTime
				clerkAvailableTime = endTime
				customer.StartTime = startTime
				customer.EndTime = endTime
				customer.ClerkID = id
				results <- customer
			}
		}(i)
	}

	var customerWg sync.WaitGroup
	customerWg.Add(len(customers))
	startTime := time.Now()

	for i := range customers {
		c := customers[i]
		go func(c Customer) {
			defer customerWg.Done()

			elapsed := time.Since(startTime)
			waitDuration := time.Duration(c.ArrivalTime)*time.Second - elapsed
			if waitDuration > 0 {
				time.Sleep(waitDuration)
			}

			numberMutex.Lock()
			currentNumber++
			numberMutex.Unlock()

			customer := &Customer{
				CustomerID:  c.CustomerID,
				ArrivalTime: c.ArrivalTime,
				ServiceTime: c.ServiceTime,
			}

			q.mutex.Lock()
			q.customers = append(q.customers, customer)
			q.cond.Signal()
			q.mutex.Unlock()
		}(c)
	}

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
