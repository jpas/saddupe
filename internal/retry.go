package internal

import (
	"time"

	"github.com/pkg/errors"
)

func Retry(times int, delay time.Duration, f func() error) error {
	var err error
	for n := 1; n <= times; n++ {
		err = f()
		if err == nil {
			break
		}
		time.Sleep(delay)
	}
	return errors.Wrapf(err, "failed after %d retries", times)
}
