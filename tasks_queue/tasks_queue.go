package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Idx  int
	Proc func() error
}

type TaskResult struct {
	Idx int
	Res error
}

func (t Task) GetResult() TaskResult {
	err := t.Proc()
	return TaskResult{Idx: t.Idx, Res: err}
}

func Worker(tasks <-chan Task, resultsChan chan<- TaskResult, wg *sync.WaitGroup) {
	for task := range tasks {
		resultsChan <- task.GetResult()
	}
	wg.Done()
}

// Run given tasks with concurrency of N
// Don't create new tasks when errlimit reached,
// Though, wait until all the existing tasks are finished
// errlimit 0 means no limit
func RunAll(tasks []func() error, N int, errlimit int) int {
	var wg = &sync.WaitGroup{}
	var tasksChan = make(chan Task)
	var resultsChan = make(chan TaskResult, N)
	var total, failures int

	// start N workers
	for i := 0; i < N; i++ {
		wg.Add(1)
		go Worker(tasksChan, resultsChan, wg)
	}

loop:
	for i := range tasks {
		select {
		case r := <-resultsChan:
			fmt.Printf("task %d result: %v\n", r.Idx, r.Res)
			total++
			if r.Res != nil {
				failures++
				if failures == errlimit {
					fmt.Println("error limit reached")
					break loop
				}
			}
			continue
		case tasksChan <- Task{i, tasks[i]}:
			fmt.Printf("Task %d starting...\n", i)
		}
	}
	close(tasksChan) // no new tasks
	wg.Wait()
	close(resultsChan) // no new results
	// drain the remaining results
	for r := range resultsChan {
		fmt.Printf("task %d result: %v\n", r.Idx, r.Res)
		total++
		if r.Res != nil {
			failures++
		}
	}
	return failures
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
