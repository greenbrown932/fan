package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// producer, worker, collector functions will be defined later

func producer(ctx context.Context, jobsCh chan<- int, numJobs int, wg *sync.WaitGroup) {
	fmt.Println("Producer: Starting")
	defer func() {
		close(jobsCh)
		fmt.Println("Producer: Closing jobs channel and exiting.")
		wg.Done()
	}()
	for i := range make([]int, numJobs) {
		jobNum := i
		select {
		case jobsCh <- jobNum:
			fmt.Println("Producer: Sent job, ", jobNum)
			time.Sleep(time.Millisecond * 500)
		case <-ctx.Done():
			fmt.Println("Producer: Context was cancelled")
			fmt.Println("Producer (error): ", ctx.Err())
			return
		}
	}
	fmt.Printf("Producer: Sent all jobs %d jobs.\n", numJobs)
}

func isChannelClosed(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// func worker(ctx context.Context, id int, jobsCh <-chan int, resultsCh chan <- string, wg *sync.WaitGroup) {
//
// }
func main() {
	fmt.Println("Starting Fan-out/Fan-in demo with Context.")
	var wg sync.WaitGroup
	numJobs := 10
	// numWorkers := 3
	pipelineTimeout := 10 * time.Second // Overall timeout for the pipeline

	ctx, cancel := context.WithTimeout(context.Background(), pipelineTimeout)
	defer cancel() // Ensure cancellation resources are cleaned up

	jobsCh := make(chan int, numJobs) // Buffer size same as numJobs can be convenient for producer
	// resultsCh := make(chan string, numJobs) // Buffer size same as numJobs

	wg.Add(1)
	go producer(ctx, jobsCh, numJobs, &wg)

	fmt.Println("Main: Simulating pipeline work or waiting for timeout...")

	isConsume := true
	jobsRecv := 0
	for isConsume {
		select {
		case job, ok := <-jobsCh:
			if ok {
				fmt.Printf("Main (temp consumer): Received job %d\n", job)
				jobsRecv++
			} else {
				fmt.Println("Main (temp consumer): Jobs channel closed.")
				isConsume = false
			}
		case <-time.After(1 * time.Second):
			fmt.Println("Main: Pipeline simulation time elapsed (longer than timeout).")
			if jobsRecv < numJobs && !isChannelClosed(ctx) { // A helper isChannelClosed would check ctx.Done
				fmt.Println("Main (temp consumer): Waited 1s for a job, but got none. Checking context or producer status.")
			}
		case <-ctx.Done():
			fmt.Printf("Main: Pipeline context canceled: %v\n", ctx.Err())
			cancel()
			isConsume = false
		}
	}
	wg.Wait()
	fmt.Println("Main: Shutting down.")
}
