package main

import (
	"context"
	"fmt"
	"time"
)

// producer, worker, collector functions will be defined later

func producer(ctx context.Context, jobsCh chan<- int, numJobs int) {
	fmt.Println("Producer: Starting")

	defer func() {
		close(jobsCh)
		fmt.Println("Producer: Closing jobs channel and exiting.")
	}()
	for i := range numJobs {
		jobNum := i
		select {
		case jobsCh <- jobNum:
			fmt.Println("Producer: Sent job, ", jobNum)
			time.Sleep(time.Millisecond * 200)
		case <-ctx.Done():
			fmt.Println("Producer: Context was cancelled")
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
func main() {
	fmt.Println("Starting Fan-out/Fan-in demo with Context.")

	// Lab: Define number of jobs and workers
	numJobs := 10
	// numWorkers := 3
	pipelineTimeout := 1 * time.Second // Overall timeout for the pipeline

	// Lab: Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), pipelineTimeout)
	defer cancel() // Ensure cancellation resources are cleaned up

	// Lab: Create channels for jobs and results
	// Use buffered channels if you anticipate slight mismatches in production/consumption rates,
	// but unbuffered is fine for understanding the core mechanics.
	jobsCh := make(chan int, numJobs) // Buffer size same as numJobs can be convenient for producer
	// resultsCh := make(chan string, numJobs) // Buffer size same as numJobs

	// var wg sync.WaitGroup // To wait for workers to finish

	// We'll launch producer, workers, and collector in the next steps.
	go producer(ctx, jobsCh, numJobs)
	// For now, let's just simulate waiting for the pipeline or timeout
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
		case <-time.After(pipelineTimeout + 1*time.Second): // Wait a bit longer than the timeout
			fmt.Println("Main: Pipeline simulation time elapsed (longer than timeout).")
			if jobsRecv < numJobs && !isChannelClosed(ctx) { // A helper isChannelClosed would check ctx.Done
				fmt.Println("Main (temp consumer): Waited 1s for a job, but got none. Checking context or producer status.")
			}
		case <-ctx.Done():
			fmt.Printf("Main: Pipeline context canceled: %v\n", ctx.Err())
		}
	}

	fmt.Println("Main: Shutting down.")
}
