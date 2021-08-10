package workersPool

import "context"

type jobID string                       //任务id
type jobType string                     //任务类型
type jobMetaData map[string]interface{} //任务metaData

//任务方法类型
type ExecFunction func(ctx context.Context, args ...interface{}) (interface{}, error)

//任务描述
type JobDescriptor struct {
	ID       jobID
	Type     jobType
	MetaData jobMetaData
}

type Job struct {
	Descriptor JobDescriptor
	Execfn     ExecFunction
	Args       interface{}
}

type Result struct {
	Value      interface{}
	Error      error
	Descriptor JobDescriptor
}

func (j *Job) excute(ctx context.Context) Result {
	result := Result{
		Descriptor: j.Descriptor,
	}
	if val, err := j.Execfn(ctx, j.Args); err != nil {
		result.Value = nil
		result.Error = err
	} else {
		result.Value = val
		result.Error = nil
	}
	return result

}
