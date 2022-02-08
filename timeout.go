package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeout = errors.New("Timeout")
)

// 引数2つ持つ関数のtimeout処理
func DoTimeoutFunc[T any, K any](parent context.Context, time time.Duration, f func(context.Context, T) (K, error), arg T) (K, error) {

	errCh := make(chan error)
	resCh := make(chan K)
	var res K

	ctx, cancel := context.WithTimeout(parent, time)
	defer cancel()

	go func() {
		r, err := f(ctx, arg)
		if err != nil {
			errCh <- err
		}
		resCh <- r
	}()

	select {
	case res = <-resCh:
		return res, nil
	case err := <-errCh:
		if err != nil {
			return res, err
		}
		return res, nil
	case <-ctx.Done():
		return res, fmt.Errorf("error: %w", ErrTimeout)
	}
}
