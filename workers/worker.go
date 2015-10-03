package workers

import (
	"fmt"
	"time"

	"github.com/teasherm/fueleconomy/global"
)

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				startTime := time.Now()
				err := work.DoWork()
				if err != nil {
					global.Logger.Println(err)
					global.Logger.Println("Task failed:", work.Target)
				} else {
					global.Logger.Println("Task succeeded:", work.Target)
					endTime := time.Now()
					global.Logger.Println(fmt.Sprintf("Time elapsed for target %s: %v",
						work.Target, endTime.Sub(startTime)))
				}
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
