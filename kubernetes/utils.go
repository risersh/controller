package kubernetes

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForResourceConditionArgs is the arguments for the WaitForResourceCondition function.
type WaitForResourceConditionArgs struct {
	// Condition is the condition to wait for such as "Ready", 123, true, false, etc.
	Condition interface{}
	// Getter is the function to get the resource.
	Getter func(ctx context.Context, name string, options metav1.GetOptions, subresources ...string) (*interface{}, error)
	// Accessor is the function to access the status of the resource.
	Accessor func(resource *interface{}) (string, bool, error)
	// Timeout is the timeout for the condition to be met.
	Timeout time.Duration
}

// WaitForResourceCondition waits for a condition to be met for each resource in the list of arguments.
//
// Arguments:
//   - args ...WaitForResourceConditionArgs: a list of arguments to wait for a condition to be met for each resource.
//
// Returns:
//   - []error: a list of errors, one for each resource.
func WaitForResourceCondition(args ...WaitForResourceConditionArgs) []error {
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
				resource, err := arg.Getter(ctx, "", metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				status, found, err := arg.Accessor(resource)
				if err != nil {
					return false, err
				}
				if !found {
					return false, nil
				}

				switch condition := arg.Condition.(type) {
				case int:
					statusInt, err := strconv.Atoi(status)
					if err != nil {
						return false, fmt.Errorf("status cannot be converted to int: %v", err)
					}
					return statusInt == condition, nil
				case bool:
					statusBool, err := strconv.ParseBool(status)
					if err != nil {
						return false, fmt.Errorf("status cannot be converted to bool: %v", err)
					}
					return statusBool == condition, nil
				case string:
					return status == condition, nil
				default:
					return false, fmt.Errorf("unsupported condition type")
				}
			})
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("condition %v not met within timeout: %v", arg.Condition, err))
				mu.Unlock()
			}
		}(arg)
	}

	wg.Wait()

	return errors
}
