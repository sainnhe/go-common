package concurrent

import (
	"sort"
	"sync"
)

/*
Run runs tasks concurrently, waits for all tasks to be executed and return ordered results.

Params:
  - concurrency int32: The maximum number of tasks that can run concurrently.
  - args []Arg: The list of arguments passed to the task handler function.
  - f func(Arg) Result: The task handler function, which receives a Arg struct and returns a Result struct.

Returns:
  - results []Result: The returned results, sorted by incoming arguments.
*/
func Run[Arg any, Result any](concurrency int32, args []Arg, f func(Arg) Result) (results []Result) {
	// Add index to the Result structure for easy sorting.
	type WrappedResult struct {
		r Result
		i int
	}

	// Start concurrent execution of tasks.
	semaphore := make(chan struct{}, concurrency)
	wrChan := make(chan WrappedResult, len(args))
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	for i, arg := range args {
		semaphore <- struct{}{}
		go func(arg Arg, i int) {
			// Use defer to avoid panic in the middle and the code is not executed.
			defer func() {
				<-semaphore
				wg.Done()
			}()
			wrChan <- WrappedResult{f(arg), i}
		}(arg, i)
	}
	go func() {
		wg.Wait()
		close(wrChan)
	}()

	// Get the execution results and sort them.
	wrs := make([]WrappedResult, 0, len(args))
	for wr := range wrChan {
		wrs = append(wrs, wr)
	}
	sort.Slice(wrs, func(i, j int) bool { return wrs[i].i < wrs[j].i })

	// Order results and return.
	results = make([]Result, 0, len(args))
	for _, wr := range wrs {
		results = append(results, wr.r)
	}
	return
}
