package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// producer, worker, collector functions will be defined later

func main() {
	fmt.Println("Starting Fan-out/Fan-in demo with Context.")

	// Lab: Define number of jobs and workers
	numJobs := 10
	numWorkers := 3
	pipelineTimeout := 5 * time.Second // Overall timeout for the pipeline

	// Lab: Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), pipelineTimeout)
	defer cancel() // Ensure cancellation resources are cleaned up

	// Lab: Create channels for jobs and results
	// Use buffered channels if you anticipate slight mismatches in production/consumption rates,
	// but unbuffered is fine for understanding the core mechanics.
	jobsCh := make(chan int, numJobs)       // Buffer size same as numJobs can be convenient for producer
	resultsCh := make(chan string, numJobs) // Buffer size same as numJobs

	var wg sync.WaitGroup // To wait for workers to finish

	// We'll launch producer, workers, and collector in the next steps.

	// For now, let's just simulate waiting for the pipeline or timeout
	fmt.Println("Main: Simulating pipeline work or waiting for timeout...")

	// Example of how main might wait (we'll refine this later)
	// This is a placeholder. In a real scenario, you'd wait for the collector
	// or specific signals, or the context to be done.
	select {
	case <-time.After(pipelineTimeout + 1*time.Second): // Wait a bit longer than the timeout
		fmt.Println("Main: Pipeline simulation time elapsed (longer than timeout).")
	case <-ctx.Done():
		fmt.Printf("Main: Pipeline context canceled: %v\n", ctx.Err())
	}

	fmt.Println("Main: Shutting down.")
}
