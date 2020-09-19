package util

import (
	"sync"
)

type WaitGroupWrapper struct {
	once sync.Once
	exitCh chan error
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func (w *WaitGroupWrapper) ExitFunc(err error){
	w.once.Do(func() {
		if err != nil {
			return
		}
		w.exitCh <- err
	})
}

