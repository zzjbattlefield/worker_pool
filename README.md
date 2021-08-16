<font size="4">

# GO-Worker-pool

## 安装
```shell
go get github.com/zzjbattlefield/worker_pool
```

## 如何使用
```go

func main() {
    //实例化pool
	pool := workersPool.New(3)
	ctx := context.Background()
    //创建worker
	go pool.Run(ctx)
    //创建任务
	go pool.GenerateJob(ctx, createJob())
	for {
		select {
		case res, ok := <-pool.Result:
			if ok {
				log.Println(res.Value)
			}
		case <-pool.Done:
			return
		}
	}
}

func createJob() []workersPool.Job {
	var testPool = make([]workersPool.Job, 0)
	for i := 0; i < 10; i++ {
		testPool = append(testPool, workersPool.Job{
			Descriptor: workersPool.JobDescriptor{Type: "any type", MetaData: nil},
			Execfn: func(ctx context.Context, args ...interface{}) (interface{}, error) {
				return args[0], nil
			},
			Args: i,
		})
	}
	return testPool
}

```
