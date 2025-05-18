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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readCustomers(filename string) ([]Customer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var customers []Customer
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 3 {
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
	return customers, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <customers_file> <num_clerks>")
		return
	}
	filename := os.Args[1]
	n, _ := strconv.Atoi(os.Args[2])

	customers, err := readCustomers(filename)
	if err != nil {
		panic(err)
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
			ticket := currentNumber
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
		fmt.Printf("%d %d %d %d\n", c.ArrivalTime, c.StartTime, c.EndTime, c.ClerkID)
	}
}
