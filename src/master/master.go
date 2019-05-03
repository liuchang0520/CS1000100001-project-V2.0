package main

import (
	"fmt"
	"os"
	// "google.golang.org/grpc"
	// "net"
	"log"
	"net/http"
	c "common"
	"net/rpc"
	"sync"

)

type Master struct { //master struct
	task string
	workerChan chan string
	closeChan chan bool //send work complete signal to shut down master
	// workerCloseChan chan bool
	client map[string]*rpc.Client //key: port running worker
	rCnt int //number of reducer
	// listener net.Listener
	input []string //input file names
}


func masterInit(task, input string, rCnt int) *Master {
	master := new(Master)
	master.workerChan = make(chan string, c.MAX_WORKER)
	master.task = task
	master.rCnt = rCnt
	master.input := getInputF(input)
	master.client = make(map[string]*rpc.Client)
	master.closeChan = make(chan bool)
	// master.workerCloseChan = make(chan bool)

	// l, err := net.Listen("tcp", c.MASTER_PORT)
	// if err != nil {
	// 	log.Fatal("master failed to listen: ", err)
	// }
	// master.listener = l

	fmt.Println("master rpc register")
	if err := rpc.Register(master); err != nil {
		log.Fatal("master rpc setup failed", err)
	}
	// fmt.Println("master rpc registered")
	rpc.HandleHTTP()
	// fmt.Println("http handler started")

	go func() {
		if err := http.ListenAndServe("localhost:" + c.MASTER_PORT, nil); err != nil {
			log.Fatal("Failed to start master process", err)
		}
	}()
	
	// fmt.Println("master setup finishes")
	return master
}

func finish(master *Master) {
	for _, c := range(master.client) {
		if err := c.Call("Worker.Close", &c.WorkArgs{}, &c.WorkerRes{}); err != nil {
			log.Fatal(err)
		}
	}
	master.closeChan <- true
}

func assignTask(master *Master, stage string) {
	/*
	map stage:
	# of input files: input files
	# of output file: # of reducer per input file (m * r)
	*/

	/*
	reduce stage:
	# of input file: m intermediate files per reducer
	# of output file: 1 per reducer (r in total)
	*/

	var taskCnt int
	if stage == c.MAP {
		if err := createInterDir(master.rCnt); err != nil {
			log.Fatal(err)
		}
		taskCnt = len(master.input)
	} else if stage == c.REDUCE {
		if err := createOutputDir(); err != nil {
			log.Fatal(err)
		}
		taskCnt = master.rCnt
	}

	var wg sync.WaitGroup //use this to wait all tasks to be finished

	for i := 0; i < taskCnt; i++ { //assign task
		wg.Add(1)
		fmt.Printf("start assigning %s task %d\n", stage, i)
		go func (taskIndex int) {
			temp := 0
			for temp < c.MAX_TEMP {
				//get an idle worker
				w <- master.workerChan
				var err error
				if err = master.client[w].Call("Worker.Work", &c.WorkArgs{}, &c.WorkerRes{}); err == nil {
					wg.Done()
					fmt.Printf("%s task %d is finished\n", stage, taskIndex)
					master.workerChan <- w
					return
				}
				fmt.Println(err.Error())
				fmt.Printf("retry assigning %s task %d\n", stage, taskIndex)
				temp++
				master.workerChan <- w
			}
			log.Fatal(fmt.SprintF("in stage %s, %d task cannot be finished\n", stage, taskIndex))
		} (i)
	}

	wg.Wait() //use this to wait all tasks to be finished
	fmt.Printf("%s stage complete\n", stage)
}

func runTask(master *Master) {
	fmt.Println("start running tasks...")
	go func() {
		assignTask(master, c.MAP)
		assignTask(master, c.REDUCE)
		aggregate() //aggregate reducer output files into a single file
		finish(master)
	}()
}

func main() {
	//args[1] indicates the map reduce task
	//args[2] indicates the input files, concatenated by '#'.e.g. 1.txt#2.txt
	//args[3] indicates the number of reducers
	if len(os.Args) != 4 {
		log.Fatal("pls specify the mapreduce task, input files, and number of reducers", " args: ", os.Args)
	}
	task := os.Args[1]
	fmt.Printf("map reduce task: %s\n", task)
	if _, ok := funcMap[task]; !ok {
		log.Fatal("invalid task, pls check the task name", fmt.Sprintf("valid task names are: %s\n", getValidTask()))
	}

	master := masterInit(task, os.Args[2], os.Args[3])
	runTask(master)
	<- master.closeChan
	// master.listener.Close()
	//send shut down signal to workers
	// <- master.workerCloseChan
	fmt.Println("master work complete")
}