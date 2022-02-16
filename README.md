# gp - Goroutine Pool

### Goroutine Pool with limited task
```go
package main

import (
	"context"
	"fmt"
	"github.com/tiancheng92/gp"
	"time"
)

func worker(input ...any) {
	if len(input) != 2 {
		fmt.Println("worker: input length error.")
		return
	}

	n1, ok := input[0].(int)
	if !ok {
		fmt.Println("worker: input.0 is not int.")
		return
	}

	n2, ok := input[1].(int)
	if !ok {
		fmt.Println("worker: input.1 is not int.")
		return
	}

	time.Sleep(time.Second)

	fmt.Printf("result: %d + %d = %d\n", n1, n2, n1+n2)
}

func limitedCountTask() {
	fmt.Println("worker pool with limited task start.")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerPool := gp.New().SetTaskCount(100).SetGoroutineCount(10).SetWorker(worker).Start(ctx)

	go func() {
		for i := 0; i < 100; i++ {
			workerPool.AddTask(i, i+1)
		}
	}()

	workerPool.Wait()

	fmt.Println("worker pool with limited task done.")
}

func main() {
	limitedCountTask()
}
```

### Goroutine Pool with infinite task
```go
package main

import (
	"context"
	"fmt"
	"github.com/tiancheng92/gp"
	"time"
)

func worker(input ...any) {
	if len(input) != 2 {
		fmt.Println("worker: input length error.")
		return
	}

	n1, ok := input[0].(int)
	if !ok {
		fmt.Println("worker: input.0 is not int.")
		return
	}

	n2, ok := input[1].(int)
	if !ok {
		fmt.Println("worker: input.1 is not int.")
		return
	}

	time.Sleep(time.Second)

	fmt.Printf("result: %d + %d = %d\n", n1, n2, n1+n2)
}

func infiniteCountTask() {
	fmt.Println("worker pool with infinite task start, continuous 20 seconds.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	workerPool := gp.New().SetGoroutineCount(10).SetWorker(worker).Start(ctx)

	go func() {
		for i := 0; ; i++ {
			workerPool.AddTask(i, i+1)
		}
	}()

	<-ctx.Done()
	fmt.Println("worker pool with infinite task done.")
}

func main() {
	infiniteCountTask()
}
```