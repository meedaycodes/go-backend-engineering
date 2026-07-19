// Package worker implements a background email worker using Go channels as a
// lightweight in-process message queue. No external broker (RabbitMQ, NATS)
// is required — the channel acts as a bounded FIFO queue between producers
// (service layer) and the single consumer goroutine.
package worker

import (
	"context"
	"log/slog"
)

// EmailJob carries the data for a single outbound email. All fields are
// strings — no transport logic here, just a data container passed over the
// channel from producer to consumer.
type EmailJob struct {
	To      string
	Subject string
	Body    string
}

// EmailWorker wraps a buffered channel of EmailJob. The buffer decouples the
// sender from the consumer — the caller does not block as long as the buffer
// has space. The channel is the queue; EmailWorker is just a named handle.
type EmailWorker struct {
	jobs chan EmailJob
}

// NewEmailWorker creates an EmailWorker with a buffered channel of the given
// size. bufferSize controls how many jobs can queue before Send blocks.
// 100 is a reasonable default for low-traffic services.
func NewEmailWorker(bufferSize int) *EmailWorker {

	emailChan := make(chan EmailJob, bufferSize)
	emailWorker := &EmailWorker{jobs: emailChan}
	return emailWorker

}

// Send enqueues a job onto the channel. Blocks if the buffer is full, so
// callers should size the buffer to handle burst traffic. In practice this
// is a fire-and-forget call — the Signup handler does not wait for delivery.
func (w *EmailWorker) Send(job EmailJob) {

	w.jobs <- job
}

// Start is the worker loop and must be run as a goroutine. It processes jobs
// one at a time. When ctx is cancelled (on server shutdown), the ctx.Done()
// case fires and the goroutine exits cleanly, draining no further jobs.
// In production the job case body would call a real email API (SendGrid, SES);
// here it logs the job for demonstration.
func (w *EmailWorker) Start(ctx context.Context) {

	for {

		select {
		case job := <-w.jobs:
			slog.Info("Processing email", "to", job.To, "subject", job.Subject)
		case <-ctx.Done():
			slog.Info("Shutting down the worker")
			return
		}
	}

}
