package utils

import (
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
