package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_ExecutesAllJobs(t *testing.T) {
	pool := NewWorkerPool(3)

	var counter atomic.Int32
	numJobs := 100

	for i := range numJobs {
		pool.Submit(func() {
			counter.Add(1)
			i++
		})
	}

	pool.Close()

	if counter.Load() != int32(numJobs) {
		t.Errorf("Expected %d jobs executed, got %d", numJobs, counter.Load())
	}
}

func TestWorkerPool_ConcurrentExecution(t *testing.T) {
	pool := NewWorkerPool(5)

	var activeWorkers atomic.Int32
	var maxConcurrent atomic.Int32

	for i := range 20 {
		pool.Submit(func() {
			current := activeWorkers.Add(1)

			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}

			time.Sleep(10 * time.Millisecond)
			activeWorkers.Add(-1)
			i++
		})
	}

	pool.Close()

	if maxConcurrent.Load() > int32(5) {
		t.Errorf("More workers executed concurrently (%d) than pool size (%d)", maxConcurrent.Load(), 5)
	}
}

func TestWorkerPool_JobOrder(t *testing.T) {
	pool := NewWorkerPool(1)
	
	var mu sync.Mutex
	var results []int

	for i := range 5 {
		jobNum := i
		pool.Submit(func() {
			mu.Lock()
			results = append(results, jobNum)
			mu.Unlock()
			i++
		})
	}

	pool.Close()

	expected := []int{0,1,2,3,4}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected results[%d] = %d, got %d", i, expected[i], v)
		}
	}
}

