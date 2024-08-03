package kubernetes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForResourceConditionArgs is the arguments for the WaitForResourceCondition function.
type WaitForResourceConditionArgs struct {
	// Evaluator is the function to access the status of the resource.
	Evaluator func() (bool, error)
	// Timeout is the timeout for the condition to be met.
	Timeout time.Duration
}

// WaitForResourceConditions waits for a condition to be met for each resource in the list of arguments.
//
// Arguments:
//   - args ...WaitForResourceConditionArgs: a list of arguments to wait for a condition to be met for each resource.
//
// Returns:
//   - []error: a list of errors, one for each resource.
func WaitForResourceConditions(args ...WaitForResourceConditionArgs) []error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	errors := make([]error, 0, len(args))

	for _, arg := range args {
		wg.Add(1)
		go func(arg WaitForResourceConditionArgs) {
			defer wg.Done()

			err := wait.PollUntilContextTimeout(ctx, time.Second, arg.Timeout, true, func(ctx context.Context) (bool, error) {
				result, err := arg.Evaluator()
				if err != nil {
					return false, err
				}
				if !result {
					return false, nil
				}

				return result, nil
			})
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("condition not met within timeout: %v", err))
				mu.Unlock()
			}
		}(arg)

	}

	wg.Wait()

	return errors
}
