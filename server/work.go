package main

import (
	"errors"
	"sync/atomic"
)

var ErrorInvalidMinRoutines = errors.New("Invalid minimum number of routines")
var ErrorWorkerPoolClosedAlready = errors.New("The WorkerPool has been already closed")

type Worker interface {
	DoWork()
}

type WorkerPool struct {
	kill     chan bool
	shutdown chan bool
	tasks    chan Worker
	routines int64 //active routines
}

func New(minRoutines int64) (*WorkerPool, error) {
	if minRoutines <= 0 {
		return nil, ErrorInvalidMinRoutines
	}

	wp := WorkerPool{
		kill:     make(chan bool),
		shutdown: make(chan bool),
		tasks:    make(chan Worker),
	}

	wp.waitShutdown()
	wp.Add(minRoutines)

	return &wp, nil
}

func (wp *WorkerPool) Run(work Worker) {
	wp.tasks <- work
}

func (wp *WorkerPool) Shutdown() error {
	if atomic.LoadInt64(&wp.routines) != 0 {
		close(wp.shutdown)
		for {
			if atomic.LoadInt64(&wp.routines) == 0 {
				return nil
			}
		}
	}
	return ErrorWorkerPoolClosedAlready
}

func (wp *WorkerPool) work() {
done:
	for {
		select {
		case t := <-wp.tasks:
			t.DoWork()
		case <-wp.kill:
			break done
		}
	}
	//shutting down
	atomic.AddInt64(&wp.routines, -1)
}

func (wp *WorkerPool) Add(nums int64) {
	if nums <= 0 {
		return
	}

	var i int64 = 1
	for ; i <= nums; i++ {
		wp.routines += 1
		go wp.work()
	}
}

func (wp *WorkerPool) waitShutdown() {
	go func() {
		for {
			select {
			case <-wp.shutdown:
				routines := int(atomic.LoadInt64(&wp.routines))
				for i := 0; i < routines; i++ {
					wp.kill <- true
				}
				return
			}
		}
	}()
}
