package utils

import "sync"

func RunFuncLConcurrently[T any](funcL []func() T) []T {
	var wg sync.WaitGroup
	results := make([]T, len(funcL))
	for i, f := range funcL {
		wg.Go(func() {
			results[i] = f()
		})
	}
	wg.Wait()
	return results
}
