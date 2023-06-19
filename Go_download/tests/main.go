package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxWorkers = 3
	NumTasks   = 10
)

type Task struct {
	ID int
}

func main() {
	// 创建任务通道和等待组
	taskChan := make(chan Task)
	var wg sync.WaitGroup

	// 启动工作池
	for i := 0; i < MaxWorkers; i++ {
		go worker(taskChan, &wg)
	}

	// 添加任务到任务通道
	for i := 0; i < NumTasks; i++ {
		task := Task{ID: i + 1}
		taskChan <- task
		wg.Add(1)
	}

	// 等待所有任务完成
	wg.Wait()
	close(taskChan)
	fmt.Println("All tasks completed.")
}

// 工作池的工作函数
func worker(taskChan <-chan Task, wg *sync.WaitGroup) {
	for task := range taskChan {
		processTask(task)
		wg.Done()
	}
}

// 模拟任务处理
func processTask(task Task) {
	fmt.Printf("Processing task %d\n", task.ID)
	time.Sleep(1 * time.Second) // 模拟任务处理时间
	fmt.Printf("Task %d completed\n", task.ID)
}
