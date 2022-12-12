package main

import (
	"context"
	"fmt"
	"github.com/tiancheng92/gp"
	"time"
)

func worker(input ...any) {
	fmt.Println("do something ...")
}

func infiniteCountTask() {
	fmt.Println("worker pool with infinite task start, continuous 10 seconds.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	workerPool := gp.New().SetGoroutineCount(1).SetWorker(worker).Start(ctx)

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			for i := 0; i < workerPool.GetGoroutineCount(); i++ {
				workerPool.AddTask()
			}
		}
	}()

	go func() {
		for i := 0; i < 9; i++ {
			time.Sleep(1 * time.Second)
			workerPool.AddGoroutineCount(1)
		}
	}()

	<-ctx.Done()
	fmt.Println("worker pool with infinite task done.")
}

func main() {
	infiniteCountTask()
}
