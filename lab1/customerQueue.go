package main

import "sync"

type CustomerQueue struct {
	customers []*Customer
	mutex     sync.Mutex
	cond      *sync.Cond
}
