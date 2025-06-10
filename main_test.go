package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

// Тест на добавление задачи при неактивных воркерах
func TestAddTaskWithoutWorkers(t *testing.T) {
	input := "add task\nend\n"
	output := captureOutput(input, func() {
		program := NewProgram()
		go program.StartProgram()
		<-program.Done
	})

	if output == "Создайте хотя бы один воркер" {
		t.Error("task added without workers")
	}
}

// Тест на добавление воркера
func TestAddWorker(t *testing.T) {
	program := NewProgram()
	input := "add worker\nend\n"

	_ = captureOutput(input, func() {
		go program.StartProgram()
		<-program.Done
	})

	if len(program.Pool.Workers) != 1 {
		t.Errorf("expected 1 worker, got %d", len(program.Pool.Workers))
	}
}

// Тест на некорректный айди
func TestRemoveWorker(t *testing.T) {
	program := NewProgram()
	program.Pool.AddWorker()

	input := strings.NewReader("remove worker\ninvalid_id\n")
	scanner := bufio.NewScanner(input)
	program.removeTaskHandler(scanner)
}

// Функция для перехвата вывода
func captureOutput(input string, f func()) string {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	f()

	wOut.Close()

	var buf bytes.Buffer
	buf.ReadFrom(rOut)
	return buf.String()
}
