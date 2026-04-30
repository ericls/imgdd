package utils

import (
	"iter"
	"sync"
)

type WorkerTask[T any] struct {
	Job     T
	Execute func(T) error
}

type workerResult[T any] struct {
	Job T
	Err error
}

func RunWorkerPool[T any](tasks []WorkerTask[T], maxWorkers int) chan workerResult[T] {
	var wg sync.WaitGroup
	jobs := make(chan WorkerTask[T], len(tasks))
	results := make(chan workerResult[T], len(tasks))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				err := task.Execute(task.Job)
				results <- workerResult[T]{Job: task.Job, Err: err}
			}
		}()
	}

	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	wg.Wait()
	close(results)

	return results
}

// RunWorkerPoolIter is like RunWorkerPool but feeds jobs from an iter.Seq2 instead of a
// pre-built slice. The jobs channel is bounded to maxWorkers*2 so at most that many items
// are held in memory at once beyond what workers are actively processing. The iterator is
// drained in a separate goroutine; if it yields an error the error is forwarded as a result
// and iteration stops.
func RunWorkerPoolIter[T any](seq iter.Seq2[T, error], execute func(T) error, maxWorkers int) chan workerResult[T] {
	var wg sync.WaitGroup
	jobs := make(chan WorkerTask[T], maxWorkers*2)
	results := make(chan workerResult[T], maxWorkers*2)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				err := task.Execute(task.Job)
				results <- workerResult[T]{Job: task.Job, Err: err}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for item, err := range seq {
			if err != nil {
				var zero T
				results <- workerResult[T]{Job: zero, Err: err}
				return
			}
			jobs <- WorkerTask[T]{Job: item, Execute: execute}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
