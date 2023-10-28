package fakenews

import (
	"context"
)

func InfiniteStream[T any](ctx context.Context, seed T, in chan T, rules ...Rule[T]) error {
	defer close(in)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fake, err := Fake[T](seed, rules...)
			if err != nil {
				return err
			}
			in <- fake
		}
	}
}
