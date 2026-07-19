package retry

import (
	"math/rand"
	"time"
)

func WithRetry(maxAttempts int, fn func() error) error {

	var err error

	for i := 0; i < maxAttempts; i++ {

		err = fn()

		if err == nil {
			return nil
		}

		wait := 100 * time.Millisecond * (1 << i)
		jitter := 0.5 + rand.Float64()

		if i < maxAttempts-1 {
			time.Sleep(time.Duration(float64(wait) * jitter))
		}
	}
	return err
}
