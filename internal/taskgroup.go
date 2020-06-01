package internal

import (
	"sync"
)

type TaskGroup struct {
	fn []func() error

	done chan struct{}
	wg   sync.WaitGroup

	errc chan error
	err  error
}

func (tg *TaskGroup) Add(f func() error) {
	tg.fn = append(tg.fn, f)
}

func (tg *TaskGroup) Done() chan struct{} {
	return tg.done
}

func (tg *TaskGroup) Start() {
	tg.done = make(chan struct{})
	tg.errc = make(chan error, len(tg.fn))

	tg.wg.Add(len(tg.fn))
	for _, f := range tg.fn {
		f := f
		go func() {
			defer tg.wg.Done()
			err := f()
			if err != nil {
				tg.Stop()
			}
			tg.errc <- err
		}()
	}
}

func (tg *TaskGroup) Stop() {
	close(tg.done)
}

func (tg *TaskGroup) Wait() {
	tg.wg.Wait()
	for err := range tg.errc {
		if err != nil {
			tg.err = err
			break
		}
	}
}

func (tg *TaskGroup) Err() error {
	tg.Wait()
	return tg.err
}

func (tg *TaskGroup) Run() error {
	tg.Start()
	tg.Wait()
	return tg.Err()
}
