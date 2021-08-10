package workersPool_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
	"workersPool"
)

func TestRun(t *testing.T) {
	pool := workersPool.New(3)
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	go pool.Run(ctx)
	go pool.GenerateJob(ctx, testJob())
	defer func() {
		time.Sleep(time.Second)
		fmt.Println("num of gorutine:", runtime.NumGoroutine())
	}()
	for {
		select {
		case val, ok := <-pool.Result:
			if ok {
				if val.Error != nil {
					// log.Println(val.Error.Error())
				} else {
					fmt.Println(val.Value)
				}

			} else {
				t.Fatalf("fail! pool.Result return not ok")
			}
		case <-pool.Done:
			return
		}
	}
}

func testJob() []workersPool.Job {
	var result = make([]workersPool.Job, 0)
	for i := 0; i < 100; i++ {
		result = append(result, workersPool.Job{
			Descriptor: workersPool.JobDescriptor{
				Type:     "anyType",
				MetaData: nil,
			},
			Execfn: func(ctx context.Context, args ...interface{}) (interface{}, error) {
				// fmt.Println(args...)
				var result int
				for _, val := range args {
					result += val.(int)
				}
				return result, nil
			},
			Args: i,
		})
	}
	return result
}
