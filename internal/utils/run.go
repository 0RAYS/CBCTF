package utils

import "sync"

func RunFuncLConcurrently[T any](funcL []func() T) []T {
	var wg sync.WaitGroup
	results := make([]T, len(funcL))
	wg.Add(len(funcL))
	for i, f := range funcL {
		go func(i int, f func() T) {
			defer wg.Done()
			results[i] = f()
		}(i, f)
	}
	wg.Wait()
	return results
}
