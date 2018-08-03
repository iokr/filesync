package util

import (
	"sync"
	"fmt"
)

type Semaphore struct {
	TotalThreads 	chan struct{}
	SemWaitGroup 	sync.WaitGroup
}

func NewSemaphore(TotalNums int) *Semaphore {
	return &Semaphore{
		TotalThreads: make(chan struct{}, TotalNums),
	}
}

func (sem *Semaphore) P() {
	sem.TotalThreads <- struct{}{}
	sem.SemWaitGroup.Add(1)
}

func (sem *Semaphore) V() {
	sem.SemWaitGroup.Done()
	<-sem.TotalThreads
}

func (sem *Semaphore) Wait() {
	sem.SemWaitGroup.Wait()
	fmt.Println("33333333333333333333333333333")
}

func (sem *Semaphore) Close() {
	close(sem.TotalThreads)
	fmt.Println("111111111111111111111111111111")
}