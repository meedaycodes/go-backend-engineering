// Package main is the entry point that ties the pipeline together:
// reads CSV data, processes records concurrently, and prints results.
// Timing is measured to demonstrate the speedup from concurrency —
// 10 records at 500ms each complete in ~500ms total, not 5000ms.
package main

import (
	"fmt"
	"time"

	"github/meedaycodes/day-03-concurrency/internal/pipeline"
)

// main reads records from CSV, runs them through the concurrent pipeline,
// and reports results. time.Now/time.Since bracket the processing to prove
// that goroutines run in parallel, not sequentially.
func main() {

	records, err := pipeline.ReadCSV("data/records.csv")

	if err != nil {
		fmt.Printf("file failed to open with error: %s.\n", err)

		return
	}

	recordCount := len(records)
	fmt.Printf("total number of records is %d.\n", recordCount)

	currentTime := time.Now()
	processedResults := pipeline.ProcessRecords(records)
	timeToProcess := time.Since(currentTime)

	fmt.Printf("it took %s seconds to process all records.\n", timeToProcess)

	for _, r := range processedResults {
		fmt.Printf("Record with ID: %d, Name: %s, status: %s.\n", r.Record.ID, r.Record.Name, r.Status)
	}
}
