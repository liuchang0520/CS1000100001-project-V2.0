package main

import (
	"fmt"
	"os"
	"log"
	// "net"
	"net/rpc"
	"net/http"
	"time"
	"errors"
	w "worker_api"
	c "common"
)

type Worker struct { //Worker Struct
	shutDownChan chan bool
	port string
	client *rpc.Client
	rCnt int //number of reducers 
	// listener net.Listener
	inputDir string
}

type RegisterArgs struct {
	Port string
}
type MasterRes struct {
	RCnt int // number of reducers we want to use
	InputDir string
}
type WorkArgs struct {
	Task string
	Stage string //map or reduce

	//for map, it is an input file; for reducer, it is a directory containing the intermediate files
	InputFile string
}
type WorkerRes struct {

}

//this rpc api is used to call worker to do a specific map/reduce task
func (worker *Worker) Work(args *WorkArgs, res *WorkerRes) error {
	fmt.Printf("worker: %s doing %s of task - %s\n", worker.port, args.Stage, args.Task)

	//get the corresponding map and reduce functions for this task
	mrFunc := c.FuncMap[args.Task]

	if args.Stage == c.MAP {
		return w.MapTask(args.Task, worker.inputDir, args.InputFile, worker.rCnt, mrFunc.MF)
	}

	if args.Stage == c.REDUCE {
		return w.ReduceTask(args.Task, args.InputFile, worker.rCnt, mrFunc.RF)
	}
	
	return errors.New(fmt.Sprint("invalid stage: %s\n", args.Stage))
}

//shutdown worker
func (worker *Worker) Close(args *WorkArgs, res *WorkerRes) error {
	worker.shutDownChan <- true
	fmt.Printf("shutdown signal received for worker: %v\n", worker.port)
	return nil
}

func workerInit(port string) *Worker {
	go func() {
		if err := http.ListenAndServe("localhost:" + port, nil); err != nil {
			log.Fatal("Failed to start worker process", err)
		}
	}()

	worker := new(Worker)
	worker.port = port
	worker.shutDownChan = make(chan bool)

	if err := rpc.Register(worker); err != nil {
		log.Fatal("worker rpc setup failed", err)
	}
	rpc.HandleHTTP()
	// fmt.Println("http handler started")
	// fmt.Println("listener started")

	workerClient, err := rpc.DialHTTP("tcp", "localhost:" + c.MASTER_PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("worker client created")
	worker.client = workerClient

	//register
	res := new(MasterRes)
	if err = workerClient.Call("Master.RegisterWorker", &RegisterArgs{Port: port}, res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("number of reduers is %d\n", res.RCnt)
	worker.rCnt = res.RCnt
	worker.inputDir = res.InputDir

	return worker
}

func main() {
	//args[1] indicates the port running worker process
	if len(os.Args) != 2 {
		log.Fatal("pls specify the worker process port number", " args: ", os.Args)
	}

	port := os.Args[1]
	fmt.Printf("worker will run at port: %s\n", port)

	worker := workerInit(port)
	<- worker.shutDownChan
	time.Sleep(5 * time.Second)
}