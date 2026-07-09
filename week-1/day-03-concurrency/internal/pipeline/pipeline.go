// Package pipeline provides functions for reading CSV data and processing
// records concurrently using the fan-out/fan-in pattern. Fan-out distributes
// work across multiple goroutines; fan-in collects results through a single channel.
// Channels are used instead of shared slices to avoid race conditions.
package pipeline

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Record represents a single row from the CSV input.
// We use a struct (not a map) because the shape is known at compile time.
// ID is int (numeric identifier), Amount is float64 (decimal money values),
// Name and Email are strings (plain text).
type Record struct {
	ID     int
	Name   string
	Email  string
	Amount float64
}

// Result pairs a processed Record with its outcome Status.
// Separate from Record because input and output are different concerns —
// the output carries processing metadata the input doesn't have.
type Result struct {
	Record Record
	Status string
}

// ReadCSV opens a CSV file, skips the header row, and converts each data row
// into a Record. Returns a slice of Records and an error.
// The file is opened with os.Open and deferred Close guarantees cleanup even
// if parsing fails partway through. strconv is used to convert string fields
// to their proper types (int, float64).
func ReadCSV(recordPath string) ([]Record, error) {

	var resultRecord []Record

	file, err := os.Open(recordPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)

	readRecord, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	slicedRecord := readRecord[1:]

	for _, row := range slicedRecord {

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, err
		}
		name := row[1]
		email := row[2]

		record := Record{ID: id, Amount: amount, Name: name, Email: email}
		resultRecord = append(resultRecord, record)
	}

	return resultRecord, nil
}

// ProcessRecord handles a single Record — simulates work with a 500ms sleep
// (standing in for a real operation like a DB write or API call) and returns
// a Result. This function processes ONE record; concurrency happens at the
// level above (ProcessRecords), where multiple goroutines each call this
// function independently.
func ProcessRecord(r Record) Result {

	time.Sleep(500 * time.Millisecond)
	fmt.Printf("The Record with the ID %d is being processed----.\n", r.ID)

	return Result{Record: r, Status: "processed"}

}

// ProcessRecords implements the fan-out/fan-in concurrency pattern:
//
// Fan-out: Each record gets its own goroutine, all running simultaneously.
// The record is passed as a parameter to the anonymous function to avoid
// the closure gotcha (all goroutines sharing the same loop variable).
//
// Channel: goroutines send Results into a channel instead of appending to
// a shared slice. This eliminates race conditions — channels are safe for
// concurrent use by design.
//
// WaitGroup: tracks how many goroutines are still running. Each goroutine
// calls Done() via defer when it finishes. A separate goroutine waits for
// the count to hit zero, then closes the channel.
//
// Fan-in: The range loop over the channel collects all results into a single
// slice. It blocks until the channel is closed, meaning all workers are done.
func ProcessRecords(rs []Record) []Result {
	results := make(chan Result)

	var wg sync.WaitGroup

	for _, r := range rs {
		wg.Add(1)

		go func(r Record) {

			defer wg.Done()
			result := ProcessRecord(r)
			results <- result

		}(r)

	}

	// Separate goroutine waits for all workers to finish, then closes the
	// channel. This must be a goroutine — if we called wg.Wait() on the main
	// goroutine here, it would block before we start reading from the channel,
	// causing a deadlock (workers can't send because nobody is receiving).
	go func() {
		wg.Wait()
		close(results)

	}()

	var processed []Result
	for r := range results {
		processed = append(processed, r)
	}

	return processed
}
