package main

import "sync"

type Customer struct {
	CustomerID  int
	ArrivalTime int
	ServiceTime int
	StartTime   int
	EndTime     int
	ClerkID     int
}

type CustomerQueue struct {
	customers []*Customer
	mutex     sync.Mutex
	cond      *sync.Cond
}
