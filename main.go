package main

import (
	"fmt"
)

func main() {
	pool := NewWorkerPool(700)

	for i := range 100000 {
		jobNum := i
		pool.Submit(func() {
			fmt.Printf("Executing job %d\n", jobNum)
		})
	}

	pool.Close()
	fmt.Println("All jobs completed")
}
