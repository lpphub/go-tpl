package ext

import (
	"context"
	"fmt"
	"testing"
)

func TestAsyncProcessor_Process(t *testing.T) {
	processor, _ := NewAsyncProcessor[int, int](func(ctx context.Context, item int) (int, error) {

		return item * 2, nil
	})

	data, _ := processor.Process(context.Background(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Printf("%v", data)
}
