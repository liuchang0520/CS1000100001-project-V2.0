package main

import (
	"fmt"
	"os"
	"log"
	// "net"
	"net/rpc"
	"net/http"
	"time"
	c "common"
)

type Worker struct { //Worker Struct
	shutDownChan chan bool
	port string
	client *rpc.Client
	rCnt int //number of reducers 
	// listener net.Listener
}

func workerInit(port string) *Worker {
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
	var res &c.MasterRes{}
	if err = workerClient.Call("Master.RegisterWorker", &c.RegisterArgs{Port: port}, res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("number of reduers is %d\n", res.RCnt)
	worker.rCnt = res.RCnt

	go func() {
		if err := http.ListenAndServe("localhost:" + port, nil); err != nil {
			log.Fatal("Failed to start worker process", err)
		}
	}()
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