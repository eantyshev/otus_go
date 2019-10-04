package main

import (
	"fmt"
	"testing"
	"time"
)

func GetWorker(waitSec int, fail bool) func() error {
	return func() error {
		d := time.Duration(waitSec) * time.Second
		time.Sleep(d)
		if fail {
			return fmt.Errorf("Task failed")
		}
		return nil
	}
}

func Scenario(t *testing.T, concurrency int, errlimit int,
	times []int, failures []bool, maxSec int, maxErr int) {
	var tasks = make([]func() error, len(times))
	for i := range times {
		tasks[i] = GetWorker(times[i], failures[i])
	}
	var result = make(chan int)
	go func() {
		result <- RunAll(tasks, concurrency, errlimit)
	}()
	select {
	case num := <-result:
		if num != maxErr {
			t.Errorf("%d tasks failed, expected %d", num, maxErr)
		}
	case <-time.After(
		time.Duration(maxSec)*time.Second +
			10*time.Millisecond):
		t.Errorf("Timeout of %d msec exceeded", maxSec)
	}
}

func TestConcurrent1(t *testing.T) {
	Scenario(t, 2, 0, []int{1, 2, 3, 4}, []bool{true, true, false, true}, 6, 3)
}

func TestErrlimit(t *testing.T) {
	Scenario(t, 2, 1, []int{1, 2, 3, 4}, []bool{true, true, false, true}, 2, 2)
}

func TestErrlimit2(t *testing.T) {
	Scenario(t, 2, 1, []int{1, 1, 1, 1}, []bool{true, false, true, true}, 2, 1)
}

func TestConcurrent2(t *testing.T) {
	Scenario(t, 10, 1, []int{4, 1, 1, 1}, []bool{true, true, true, true}, 4, 4)
}
