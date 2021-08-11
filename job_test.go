package workersPool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	descriptor = JobDescriptor{
		ID:   jobID("1"),
		Type: jobType("anyType"),
		MetaData: jobMetaData{
			"foo": "foo",
			"bar": "bar",
		},
	}
	execFn = func(ctx context.Context, args ...interface{}) (interface{}, error) {
		argVal, ok := args[0].(int)
		if !ok {
			return nil, errDefault
		}
		return argVal * 2, nil
	}
)

func TestExecute(t *testing.T) {
	ctx := context.TODO()
	type field struct {
		descriptor JobDescriptor
		execFn     ExecFunction
		args       interface{}
	}
	test := []struct {
		name   string
		fields field
		want   Result
	}{
		{
			name: "success",
			fields: field{
				descriptor: descriptor,
				execFn:     execFn,
				args:       10,
			},
			want: Result{
				Value:      20,
				Descriptor: descriptor,
			},
		},
		{
			name: "fial",
			fields: field{
				descriptor: descriptor,
				execFn:     execFn,
				args:       "10",
			},
			want: Result{
				Error:      errDefault,
				Descriptor: descriptor,
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			j := Job{
				Descriptor: tt.fields.descriptor,
				Execfn:     tt.fields.execFn,
				Args:       tt.fields.args,
			}
			got := j.execute(ctx)
			if tt.want.Error != nil {
				if !reflect.DeepEqual(got.Error, tt.want.Error) {
					t.Errorf("got error:%v want error:%v", got.Error, tt.want.Error)
				}
				// return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
