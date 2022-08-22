package gp

import (
	"context"
	"sync"
)

// goroutinePool is a pool of goroutines.
type goroutinePool struct {
	goroutineCount  int          // goroutineCount is the number of goroutines in the pool.
	inputChannel    chan []any   // inputChannel is the channel that the goroutines will read from.
	inputChannelCap int          // inputChannelCap is the capacity of the input channel.
	worker          func(...any) // worker is the function that the goroutines will execute.
}

// New creates a goroutine pool.
func New() *goroutinePool {
	return &goroutinePool{}
}

// SetGoroutineCount sets the goroutine count.
func (g *goroutinePool) SetGoroutineCount(count int) *goroutinePool {
	g.goroutineCount = count
	return g
}

// SetSizeOfWorkerQueue sets the size of the worker queue.
func (g *goroutinePool) SetSizeOfWorkerQueue(count int) *goroutinePool {
	g.inputChannelCap = count
	return g
}

// SetWorker sets the worker function.
func (g *goroutinePool) SetWorker(worker func(...any)) *goroutinePool {
	g.worker = worker
	return g
}

// check the goroutine pool for errors.
func (g *goroutinePool) check() {
	if g.worker == nil {
		panic("worker is nil")
	}

	if g.goroutineCount < 1 {
		panic("goroutine count must be greater than 0")
	}
}

// SetTaskCount sets the task count.
func (g *goroutinePool) SetTaskCount(count int) *goroutinePoolLimitedTaskCount {
	gl := &goroutinePoolLimitedTaskCount{
		goroutinePool: *g,
	}
	gl.wg.Add(count)
	return gl
}

// AddTask adds a task to the goroutine pool.
func (g *goroutinePool) AddTask(input ...any) {
	g.inputChannel <- input
}

// Start starts the goroutine pool.
func (g *goroutinePool) Start(ctx context.Context) *goroutinePool {
	g.check()
	g.inputChannel = make(chan []any, g.inputChannelCap)
	once := &sync.Once{}
	for i := 0; i < g.goroutineCount; i++ {
		go func() {
			for {
				select {
				case inputValue, ok := <-g.inputChannel:
					if !ok {
						return
					}
					g.worker(inputValue...)
				case <-ctx.Done():
					once.Do(func() {
						close(g.inputChannel)
					})
					return
				}
			}
		}()
	}
	return g
}

// goroutinePoolLimitedTaskCount is a pool of goroutines with limit task.
type goroutinePoolLimitedTaskCount struct {
	goroutinePool
	wg sync.WaitGroup
}

// SetGoroutineCount sets the goroutine count.
func (g *goroutinePoolLimitedTaskCount) SetGoroutineCount(count int) *goroutinePoolLimitedTaskCount {
	g.goroutineCount = count
	return g
}

// SetSizeOfWorkerQueue sets the size of the worker queue.
func (g *goroutinePoolLimitedTaskCount) SetSizeOfWorkerQueue(count int) *goroutinePoolLimitedTaskCount {
	g.inputChannelCap = count
	return g
}

// SetWorker sets the worker function.
func (g *goroutinePoolLimitedTaskCount) SetWorker(worker func(...any)) *goroutinePoolLimitedTaskCount {
	g.worker = worker
	return g
}

// check the goroutine pool for errors.
func (g *goroutinePoolLimitedTaskCount) check() {
	if g.worker == nil {
		panic("worker is nil")
	}

	if g.goroutineCount < 1 {
		panic("goroutine count must be greater than 0")
	}
}

// AddTask adds a task to the goroutine pool.
func (g *goroutinePoolLimitedTaskCount) AddTask(input ...any) {
	g.inputChannel <- input
}

// Start starts the goroutine pool.
func (g *goroutinePoolLimitedTaskCount) Start(ctx context.Context) *goroutinePoolLimitedTaskCount {
	g.check()
	g.inputChannel = make(chan []any, g.inputChannelCap)
	once := &sync.Once{}
	for i := 0; i < g.goroutineCount; i++ {
		go func() {
			for {
				select {
				case inputValue, ok := <-g.inputChannel:
					if !ok {
						return
					}
					g.worker(inputValue...)
					g.wg.Done()
				case <-ctx.Done():
					once.Do(func() {
						close(g.inputChannel)
					})
					return
				}
			}
		}()
	}
	return g
}

// Wait waits for all tasks to be completed.
func (g *goroutinePoolLimitedTaskCount) Wait() {
	g.wg.Wait()
}
