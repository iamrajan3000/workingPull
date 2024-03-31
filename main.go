package main

import (
	"fmt"
	"sync"
)

func processData(id int, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	// Simulating some processing time
	result := fmt.Sprintf("Processed data from ID %d", id)
	ch <- result
}

func processBatch(startID, batchSize int, ch chan<- string) {
	var wg sync.WaitGroup                       // Create a new WaitGroup for each batch
	semaphore := make(chan struct{}, batchSize) // Semaphore to limit concurrent goroutines

	for i := startID; i < startID+batchSize; i++ {
		semaphore <- struct{}{} // Acquire semaphore
		wg.Add(1)
		go func(id int) {
			defer func() {
				<-semaphore // Release semaphore
			}()
			processData(id, &wg, ch)
		}(i)
	}

	wg.Wait() // Wait for all goroutines in this batch to complete
	close(ch) // Close the channel after all processing is done for this batch
}

func main() {
	// Total number of iterations
	totalIterations := 1000
	// Batch size
	batchSize := 5

	for i := 1; i <= totalIterations; i += batchSize {
		// Channel for communication between goroutines
		ch := make(chan string)

		// Start a goroutine to process a batch of data
		go processBatch(i, batchSize, ch)

		// Receive and print results from the channel for each batch
		for result := range ch {
			fmt.Println(result)
		}

		fmt.Println("Batch Completed") // Print a separator between batches
	}
}
