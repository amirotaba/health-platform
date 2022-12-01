package wpool

import (
	"context"
)

type JobID string
type jobType string
type jobMetadata map[string]interface{}

type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)
type ExecutionFn2 func(ctx context.Context, args interface{}) error

type JobDescriptor struct {
	ID       JobID
	JType    jobType
	Metadata map[string]interface{}
}

type Result struct {
	Value      interface{}
	Err        error
	Descriptor JobDescriptor
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn2
	Args       interface{}
}

func (j Job) execute(ctx context.Context) Result {
	err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Descriptor: j.Descriptor,
	}
}
