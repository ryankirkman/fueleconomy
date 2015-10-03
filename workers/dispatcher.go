package workers

import (
	"github.com/teasherm/fueleconomy/global"
)

var WorkQueue = make(chan WorkRequest, 1000)

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	for i := 0; i < nworkers; i++ {
		global.Logger.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				global.Logger.Println("Received work request for:", work.Target)
				go func() {
					worker := <-WorkerQueue
					global.Logger.Println("Dispatching work request for:", work.Target)
					worker <- work
				}()
			}
		}
	}()
}
