package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tousart/workerpool/workerpool"
)

type Program struct {
	Pool *workerpool.WorkerPool
	Done chan struct{}
}

func NewProgram() *Program {
	return &Program{
		Pool: workerpool.NewWorkerPool(),
		Done: make(chan struct{}),
	}
}

func (p *Program) StartProgram() {
	defer close(p.Done)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			fmt.Println("Конец работы")
			return
		}

		cmd := scanner.Text()

		if cmd == "add worker" {

			p.Pool.AddWorker()

		} else if cmd == "add task" {

			p.addTaskHandler(scanner)

		} else if cmd == "remove worker" {

			p.removeTaskHandler(scanner)

		} else if cmd == "end" {

			close(p.Pool.TasksQueue)
			p.Pool.SyncGroup.Wait()
			return

		} else {
			fmt.Println("Неправильная команда")
		}
	}
}

func (p *Program) addTaskHandler(scanner *bufio.Scanner) {
	if len(p.Pool.Workers) == 0 {
		fmt.Println("Создайте хотя бы один воркер")
		return
	}

	fmt.Print("Введите задачу: ")
	if !scanner.Scan() {
		fmt.Println("Конец работы")
		return
	}

	task := strings.TrimSpace(scanner.Text())
	if len(task) == 0 {
		fmt.Println("Задача не введена")
		return
	}

	p.Pool.AddTask(task)
}

func (p *Program) removeTaskHandler(scanner *bufio.Scanner) {
	fmt.Print("Введите айди воркера: ")
	if !scanner.Scan() {
		fmt.Println("Конец работы")
		return
	}
	idText := scanner.Text()
	id, err := strconv.Atoi(idText)
	if err != nil {
		fmt.Printf("Неправильный айди: %v\n", err)
		return
	}
	p.Pool.RemoveWorker(id)
}

func main() {
	program := NewProgram()

	go program.StartProgram()

	<-program.Done
}
