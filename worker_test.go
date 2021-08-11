package workersPool

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	pool := New(3)
	ctx := context.Background()
	go pool.Run(ctx)
	go pool.GenerateJob(ctx, testJob())
	defer func() {
		time.Sleep(time.Microsecond)
		if runtime.NumGoroutine() != 2 {
			t.Fatalf("存在goroutine泄露 退出时剩余:%d", runtime.NumGoroutine())
		}
	}()
	for {
		select {
		case val, ok := <-pool.Result:
			if ok {
				if val.Error != nil {
					log.Println(val.Error.Error())
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

func TestRun_Cancel(t *testing.T) {
	pool := New(5)
	ctx, cancel := context.WithCancel(context.Background())
	go pool.Run(ctx)
	cancel()
	defer func() {
		time.Sleep(time.Microsecond)
		if runtime.NumGoroutine() != 2 {
			t.Fatalf("存在goroutine泄露 退出时剩余:%d", runtime.NumGoroutine())
		}
	}()
	for {
		select {
		case r, ok := <-pool.Result:
			if !ok {
				fmt.Println("all result recive")
			}
			if r.Error != nil && r.Error != context.Canceled {
				t.Fatalf("want error:%v,got error:%v", context.Canceled, r.Error)
			}
		case <-pool.Done:
			return
		}
	}
}

func testJob() []Job {
	var result = make([]Job, 0)
	for i := 0; i < 100; i++ {
		result = append(result, Job{
			Descriptor: JobDescriptor{
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
