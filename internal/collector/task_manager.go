package collector

import (
	"context"
	"sync"
)

type TaskManager struct {
	taskQueue chan ForecastTask
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// Initializes the TaskManager and the context
func NewTaskManager(workerCount int) *TaskManager {
	ctx, cancel := context.WithCancel(context.Background()) // Start context with cancel function
	return &TaskManager{
		taskQueue: make(chan ForecastTask, workerCount),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Initialize the worker's goroutines
func (tm *TaskManager) StartWorkers(workerCount int) {
	for range workerCount {
		tm.wg.Add(1)
		go func() {
			defer tm.wg.Done()
			for {
				select {
				case <-tm.ctx.Done(): // Checks if the context is canceled, signaling the worker to stop
					return
				case task := <-tm.taskQueue: // Checks if there is a task in the taskQueue channel
					result, err := fetchWeatherForecast(task.Api, task.Lat, task.Lon, task.Day) // Processes the task by fetching the weather forecast
					if err != nil {
						task.Err <- err // Sends the error back through the task's error channel.
					} else {
						task.Result <- result // Sends the result back through the task's result channel.
					}
				}
			}
		}()
	}
}

// Stops all workers and waits for them to finish
func (tm *TaskManager) StopWorkers() {
	tm.cancel()
	tm.wg.Wait()
}

// Adds a task to the task queue
func (tm *TaskManager) AddTask(task ForecastTask) {
	tm.taskQueue <- task
}
