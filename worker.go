package main

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	numOfWorkers int
	jobs         chan func()
	wg           sync.WaitGroup
}

func NewWorkerPool(numOfWorkers int) *WorkerPool {
	pool := &WorkerPool{
		numOfWorkers: numOfWorkers,
		jobs: make(chan func()),
	}

	for i := range numOfWorkers {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for job := range p.jobs {
		fmt.Printf("Worker %d: starting job\n", id)
		job()
		fmt.Printf("Worker %d: finished job\n", id)
	} 
}

func (p *WorkerPool) Submit(job func()) {
	p.jobs <- job
}

func (p *WorkerPool) Close() {
	close(p.jobs)
	p.wg.Wait()
}
