package work_queue

import ()

type Worker interface {
	Run() interface{}
}

type WorkQueue struct {
	Jobs         chan Worker
	Results      chan interface{}
	StopRequests chan int
	NumWorkers   uint
}

// Create a new work queue capable of doing nWorkers simultaneous tasks, expecting to queue maxJobs tasks.
func Create(nWorkers uint, maxJobs uint) *WorkQueue {
	q := WorkQueue{
		Jobs:         make(chan Worker, maxJobs),
		Results:      make(chan interface{}),
		StopRequests: make(chan int, nWorkers),
		NumWorkers:   nWorkers,
	}
	for i := 1; i <= int(nWorkers); i++ {
		go q.worker()
	}
	return &q
}

// A worker goroutine that processes tasks from .Jobs unless .StopRequests has a message saying to halt now.
func (queue WorkQueue) worker() {
	running := true
	// Run tasks from the queue, unless we have been asked to stop.

	for running {
		select {
		case j := <-queue.Jobs:
			res := j.Run()
			queue.Results <- res
			continue
		default:
		}
		select {
		case j := <-queue.Jobs:
			res := j.Run()
			queue.Results <- res
			continue
		case <-queue.StopRequests:
			running = false
		}
		break
	}
}

func (queue WorkQueue) Enqueue(work Worker) {
	queue.Jobs <- work
}

func (queue WorkQueue) Shutdown() {
	// When the queue's .Shutdown method is called, .NumWorkers messages should be sent on the .StopRequests channel.
	for i := 1; i <= int(queue.NumWorkers); i++ {
		queue.StopRequests <- 1
	}
}
