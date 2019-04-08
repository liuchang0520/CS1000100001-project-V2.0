package main

import (
	"fmt"
	"os"
	"log"
	"net/rpc"
	"net/http"
	c "common"
)

type Worker struct { //Worker Struct
	shutDownChan chan bool
}
//implement worker's rpc api
//this rpc api is used to call worker to do a specific map/reduce task
func (worker *Worker) Work(args *c.WorkArgs, res *c.WorkerRes) error {
	return nil
}


func workerInit(port string) {
	worker := new(Worker)

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
	fmt.Println("client created")

	//register
	if err = workerClient.Call("Master.RegisterWorker", &c.RegisterArgs{Port: port}, &c.MasterRes{}); err != nil {
		log.Fatal(err)
	}

	////TODO: keep listener in worker struct
	go func() {if err := http.ListenAndServe("localhost:" + port, nil); err != nil {
			log.Fatal("Failed to start worker process", err)
	}}()
	<- worker.shutDownChan
}

func main() {
	//args[1] indicates the port running worker process
	if len(os.Args) != 2 {
		log.Fatal("pls specify the worker process port number", " args: ", os.Args)
	}

	port := os.Args[1]
	fmt.Printf("worker will run at port: %s\n", port)

	workerInit(port)
}