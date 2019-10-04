package main

import (
	"fmt"
	"sync"
	"time"
)

type TaskResult struct {
	itask  int
	result error
}

// Run given tasks with concurrency of N
// Don't create new tasks when errlimit reached,
// Though, wait until all the existing tasks are finished
// errlimit 0 means no limit
func RunAll(tasks []func() error, N int, errlimit int) int {
	var wg = &sync.WaitGroup{}
	var runPool = make(chan struct{}, N)
	var results = make(chan TaskResult, N)
	var term = make(chan struct{}, 1)
	var errcount = make(chan int)
	go func() {
		var cnt = 0
		for r := range results {
			fmt.Printf("task %d result: %v\n", r.itask, r.result)
			if r.result != nil {
				cnt++
				if cnt == errlimit {
					term <- struct{}{}
				}
			}
			<-runPool
		}
		errcount <- cnt
	}()
loop:
	for i, task := range tasks {
		select {
		case <-term:
			fmt.Println("error limit reached, no new tasks")
			break loop
		case runPool <- struct{}{}:
			fmt.Printf("Task %d starting...\n", i)
			wg.Add(1)
			go func(task func() error, itask int, wg *sync.WaitGroup) {
				defer wg.Done()
				results <- TaskResult{itask: itask, result: task()}
			}(task, i, wg)
		}
	}
	close(runPool) // no new tasks
	wg.Wait()
	close(results) // no new results
	return <-errcount
}

func main() {
	var task0 = func() error { time.Sleep(time.Second); return fmt.Errorf("xcxccxcx") }
	var task1 = func() error { time.Sleep(2 * time.Second); return fmt.Errorf("xcxccxcx") }
	var task2 = func() error { time.Sleep(3 * time.Second); return fmt.Errorf("xcxccxcx") }
	var task3 = func() error { time.Sleep(4 * time.Second); return fmt.Errorf("xcxccxcx") }
	var tasks = []func() error{task0, task1, task2, task3}
	fails := RunAll(tasks, 4, 1)
	fmt.Println("RunAll: ", fails)
}
