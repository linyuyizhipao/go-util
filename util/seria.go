package util

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)



func try(fn func(), cleaner func()) (err error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		_, file, line, _ := runtime.Caller(2)
		if rErr := recover(); rErr != nil {
			if _, ok := rErr.(error); ok {
				err = errors.New(fmt.Sprintf("%s:%d,err:%v", file, line,rErr.(error)))
			} else {
				err = fmt.Errorf("%+v", rErr)
			}
		}
	}()

	fn()

	return nil
}


// Parallel 并发执行
func Parallel(fns ...func()) func() {
	var wg sync.WaitGroup
	return func() {
		wg.Add(len(fns))
		for _, fn := range fns {
			go func() {
				if err :=try(fn, wg.Done);err!=nil{
					return
				}
			}()
		}
		wg.Wait()
	}
}

// Serial 串行
func Serial(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}
