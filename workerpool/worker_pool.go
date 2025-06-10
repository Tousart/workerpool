package workerpool

import (
	"fmt"
	"sync"

	"github.com/tousart/workerpool/worker"
)

type WorkerPool struct {
	Workers    map[int]*worker.Worker
	SyncGroup  sync.WaitGroup
	TasksQueue chan string
	WorkerID   int
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		Workers:    make(map[int]*worker.Worker),
		TasksQueue: make(chan string),
	}
}

func (p *WorkerPool) AddWorker() {
	id := p.WorkerID
	wrkr := worker.NewWorker(id, p.TasksQueue)
	p.Workers[id] = wrkr
	p.WorkerID++
	p.SyncGroup.Add(1)

	go wrkr.StartWorker(&p.SyncGroup)
	fmt.Printf("воркер %d начал работу\n", id)
}

func (p *WorkerPool) RemoveWorker(id int) {
	worker, exists := p.Workers[id]
	if !exists {
		fmt.Printf("воркера %d не существует\n", id)
		return
	}

	worker.StopWorker()
	delete(p.Workers, id)
}

func (p *WorkerPool) AddTask(task string) {
	p.TasksQueue <- task
}
