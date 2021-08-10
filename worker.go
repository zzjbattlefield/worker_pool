package workersPool

import (
	"context"
	"log"
	"sync"
)

type workersPool struct {
	workersCount int
	jobs         chan Job
	Result       chan Result
	Done         chan struct{}
}

func New(wCount int) *workersPool {
	return &workersPool{
		workersCount: wCount,
		jobs:         make(chan Job),
		Result:       make(chan Result),
		Done:         make(chan struct{}),
	}
}

func (p *workersPool) Run(ctx context.Context) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < p.workersCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg)
	}
	wg.Wait()
	//所有任务已经执行完毕
	p.Done <- struct{}{}
	close(p.Result)
}

func (p *workersPool) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			//context取消
			log.Printf("context cancel:%v", ctx.Err().Error())
			p.Result <- Result{
				Error: ctx.Err(),
			}
			return
		case job, ok := <-p.jobs:
			if !ok {
				//任务队列被关闭 worker退出
				return
			} else {
				result := job.excute(ctx)
				p.Result <- result
			}
		}
	}
}

//任务生成器
func (p *workersPool) GenerateJob(ctx context.Context, jobs []Job) {
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			close(p.jobs)
			return
		default:
			p.jobs <- job
		}
	}
	//分配完所有任务后关闭jobChannel worker做完所有任务后自动退出
	close(p.jobs)
}
