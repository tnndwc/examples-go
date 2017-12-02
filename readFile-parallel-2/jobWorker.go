package main

import (
	"bufio"
	"os"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
	"container/list"
	"time"
)

const (
	queueCapacity = 50
)

func ReadFile(path string) {
	fmt.Println(path)
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var bf []byte
	capVal := 1 * 1 * 1024
	bf = make([]byte, capVal, capVal)
	scanner.Buffer(bf, bufio.MaxScanTokenSize*10)
	for scanner.Scan() {
		scanner.Bytes()
	}
}

type Job struct {
	filePath string
	wait     *sync.WaitGroup
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				ReadFile(job.filePath)
				job.wait.Done()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	w.quit <- true
	/*go func() {
		w.quit <- true
	}()*/
}

type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run() *list.List {
	workerList := list.New()
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		workerList.PushBack(worker)
		worker.Start()
	}
	go d.dispatch()
	return workerList
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}

var JobQueue chan Job = make(chan Job, queueCapacity)

func main() {
	start := time.Now()

	//---------

	var rootPath string
	flag.StringVar(&rootPath, "path", "", "the path to read files")
	flag.Parse()

	if len(rootPath) <= 0 {
		fmt.Println("please enter the path")
		return
	}

	fmt.Println("path: " + rootPath)

	dispatcher := NewDispatcher(4)
	workList := dispatcher.Run()

	path := rootPath
	files, _ := ioutil.ReadDir(rootPath)

	appWg := new(sync.WaitGroup)

	allFileList := list.New()

	for _, f := range files {
		if !f.IsDir() {
			allFileList.PushBack(path + "/" + f.Name())
			appWg.Add(1)
		}
	}

	go func() {
		var fp string
		for f := allFileList.Front(); f != nil; f = f.Next() {
			fp = f.Value.(string)

			job := Job{filePath: fp, wait: appWg}
			JobQueue <- job
		}
	}()

	appWg.Wait()

	//--end----

	end := time.Since(start)
	fmt.Println("App: ", end)

	for workerEle := workList.Front(); workerEle != nil; workerEle = workerEle.Next() {
		workerEle.Value.(Worker).Stop()
	}

}
