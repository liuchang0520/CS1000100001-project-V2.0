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
)

type Master struct { //master struct
	task string
	workerChan chan string
	closeChan chan bool //send work complete signal to shut down master
	// workerCloseChan chan bool
	client []*rpc.Client
	// listener net.Listener
}


func masterInit(task string) *Master {
	master := new(Master)
	master.workerChan = make(chan string, c.MAX_WORKER)
	master.task = task
	master.client = make([]*rpc.Client, MAX_WORKER)
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

func runTask(master *Master) {
	fmt.Println("start running tasks...")
	go func() {
		assignMapTask()
		assignReduceTask()
		finish(master)
	}()
}

func main() {
	//args[1] indicates the map reduce task
	if len(os.Args) != 2 {
		log.Fatal("pls specify the mapreduce task", " args: ", os.Args)
	}
	task := os.Args[1]
	fmt.Printf("map reduce task: %s\n", task)

	master := masterInit(task)
	runTask(master)
	<- master.closeChan
	// master.listener.Close()
	//send shut down signal to workers
	// <- master.workerCloseChan
	fmt.Println("master work complete")
}