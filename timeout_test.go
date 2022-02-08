package timeout

import (
	"context"
	"errors"
	"testing"
	"time"
)

type testCase[T any, K any] struct {
	name   string
	time   time.Duration
	arg    T
	f      func(context.Context, T) (K, error)
	expect struct {
		res K
		err error
	}
}

func TestDotimeoutFunc(t *testing.T) {

	testCases := []testCase[int, int]{
		{
			name: "success case",
			time: time.Second * 2,
			f: func(ctx context.Context, b int) (int, error) {
				return b, nil
			},
			arg: 8,
			expect: struct {
				res int
				err error
			}{
				res: 8,
				err: nil,
			},
		},
		{
			name: "timeout case",
			time: time.Second,
			f: func(ctx context.Context, a int) (int, error) {
				time.Sleep(time.Second * 2)
				return a, nil
			},
			arg: 8,
			expect: struct {
				res int
				err error
			}{
				res: 8,
				err: ErrTimeout,
			},
		},
	}
	for _, v := range testCases {
		res, err := DoTimeoutFunc(context.Background(), v.time, v.f, v.arg)
		if !errors.Is(err, v.expect.err) {
			t.Errorf("test %s: expect err: %+v, got %+v\n", v.name, v.expect.err, err)
			continue
		}
		if res != v.expect.res {
			t.Errorf("test %s: expect res: %+v, got %+v\n", v.name, v.expect.res, res)
		}
	}
}
