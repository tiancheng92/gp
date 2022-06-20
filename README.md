
# 项目简介

* gf库提供了快速构建golang协程池方法，支持无限任务模型与有限任务模型。
* 项目地址：[tiancheng92/gp](https://github.com/tiancheng92/gp)

<!--more-->

# 使用方法

```shell
go get -u github.com/tiancheng92/gp
```

## 有限任务模型

```golang
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

## 无限任务模型

```golang
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

# gp库方法详解

| 函数签名                                                                                                   | 描述                        |
|:-------------------------------------------------------------------------------------------------------|---------------------------|
| func New() *goroutinePool                                                                              | 创建goroutinePool对象（无限任务模型） |
| func (g *goroutinePool) SetGoroutineCount(count int) *goroutinePool                                    | 设置Goroutine数量（无限任务模型）     |
| func (g *goroutinePool) SetSizeOfWorkerQueue(count int) *goroutinePool                                 | 设置队列容量（无限任务模型）            |
| func (g *goroutinePool) SetWorker(worker func(...any)) *goroutinePool                                  | 设置worker方法（无限任务模型）        |
| func (g *goroutinePool) Start(ctx context.Context) *goroutinePool                                      | 启动协程池（无限任务模型）             |
| func (g *goroutinePool) AddTask(input ...any)                                                          | 添加任务（无限任务模型）              |
| func (g *goroutinePool) SetTaskCount(count int) *goroutinePoolLimitedTaskCount                         | 设置任务数量（会把无限任务模型转化为有限任务模型） |
| func (g *goroutinePoolLimitedTaskCount) SetGoroutineCount(count int) *goroutinePoolLimitedTaskCount    | 创建goroutinePool对象（有限任务模型） |
| func (g *goroutinePoolLimitedTaskCount) SetSizeOfWorkerQueue(count int) *goroutinePoolLimitedTaskCount | 设置队列容量（有限任务模型）            |
| func (g *goroutinePoolLimitedTaskCount) SetWorker(worker func(...any)) *goroutinePoolLimitedTaskCount  | 设置worker方法（有限任务模型）        |
| func (g *goroutinePoolLimitedTaskCount) AddTask(input ...any)                                          | 添加任务（有限任务模型）              |
| func (g *goroutinePoolLimitedTaskCount) Start(ctx context.Context) *goroutinePoolLimitedTaskCount      | 启动协程池（有限任务模型）             |
| func (g *goroutinePoolLimitedTaskCount) Wait()                                                         | 等待任务完成 （阻塞主进程）            |
