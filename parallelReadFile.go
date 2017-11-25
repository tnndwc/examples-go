package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"flag"
	"sync"
	"container/list"
	"path"
	"runtime"
)

const zero = 0

func startProcess(path string) {
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var bf []byte
	capVal := 1 * 1 * 1024
	bf = make([]byte, capVal, capVal)
	//scanner.Buffer(bf, bufio.MaxScanTokenSize*10)
	scanner.Buffer(bf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		//scanner.Bytes()
	}
}

func job(appendJob chan int, filePath string, wg *sync.WaitGroup) {
	defer func() {
		<-appendJob
		wg.Done()
	}()
	fmt.Println(" read file : " + filePath)
	startProcess(filePath)
}

func main() {
	var rootPath string
	flag.StringVar(&rootPath, "path", "", "the path to read files")
	flag.Parse()

	if len(rootPath) <= 0 {
		fmt.Println("please enter the path")
		return
	}

	fmt.Println("path: " + rootPath)

	files, _ := ioutil.ReadDir(rootPath)

	jobQueueSize := runtime.NumCPU()

	jobQueue := make(chan int, jobQueueSize)

	appWg := new(sync.WaitGroup)

	fileList := list.New()

	for _, f := range files {
		if !f.IsDir() {
			fileList.PushBack(path.Join(rootPath, f.Name()))
			appWg.Add(1)
		}
	}

	start := time.Now()
	go func() {
		for f := fileList.Front(); f != nil; f = f.Next() {
			jobQueue <- zero
			var fp string
			fp = f.Value.(string)
			go job(jobQueue, fp, appWg)
		}
	}()

	appWg.Wait()
	end := time.Since(start)
	fmt.Println("App: ", end)
}
