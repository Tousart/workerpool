package worker

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id    int
	tasks chan string
	stop  chan bool
}

func NewWorker(id int, tasks chan string) *Worker {
	return &Worker{
		id:    id,
		tasks: tasks,
		stop:  make(chan bool, 1),
	}
}

func (w *Worker) StartWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-w.stop:
			fmt.Printf("сигнал остановки - воркер %d завершил работу\n", w.id)
			return

		case task, ok := <-w.tasks:
			if !ok {
				fmt.Printf("канал закрыт - воркер %d завершил работу\n", w.id)
				return
			}
			fmt.Printf("воркер %d обрабатывает задачу: %s\n", w.id, task)
			time.Sleep(3 * time.Second) // Бурная деятельность
			fmt.Printf("воркер %d выполнил задачу: %s\n", w.id, task)
		}
	}
}

func (w *Worker) StopWorker() {
	w.stop <- true
}
