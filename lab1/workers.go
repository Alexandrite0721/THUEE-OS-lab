package main

import (
	"sync"
	"time"
)

func startClerkWorkers(n int, q *CustomerQueue, results chan *Customer) *sync.WaitGroup {
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
	return &clerkWg
}

func startCustomerWorkers(customers []Customer, q *CustomerQueue) *sync.WaitGroup {
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
	return &customerWg
}