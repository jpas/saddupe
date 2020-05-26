package main

import "sync"

func waitAll(fn ...func(<-chan struct{}) error) error {
	stop := make(chan struct{})
	result := make(chan error, len(fn))

	if len(fn) < 1 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(fn))

	for _, f := range fn {
		f := f
		go func() {
			defer wg.Done()
			result <- f(stop)
		}()
	}

	// defer to let others finish in the background
	defer wg.Wait()
	defer close(stop)

	return <-result
}
