package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"net/rpc"
)

type Master struct { //master struct
	task string
	workerChan chan string
}

//implement master's rpc api
//this rpc api is used to register workers
func (master *Master) RegisterWorker(args *RegisterArgs, res *MasterRes) error {
	fmt.Println(args.Port, " registering")
	master.workerChan <- args.Port
	fmt.Printf("worker at port: %s registered\n", args.Port)
	return nil
}

func masterInit(task string) {
	master := new(Master)
	master.workerChan = make(chan string, MAX_WORKER)
	master.task = task

	fmt.Println("master rpc register")
	if err := rpc.Register(master); err != nil {
		log.Fatal("master rpc setup failed", err)
	}
	// fmt.Println("master rpc registered")
	rpc.HandleHTTP()
	// fmt.Println("http handler started")


	//TODO: keep listener in master struct: l, e := net.Listen("tcp", ":1234")
	if err := http.ListenAndServe("localhost:" + MASTER_PORT, nil); err != nil {
		log.Fatal("Failed to start master process", err)
	}
	// fmt.Println("master setup finishes")
}

func main() {
	//args[1] indicates the map reduce task
	if len(os.Args) != 2 {
		log.Fatal("pls specify the mapreduce task", " args: ", os.Args)
	}
	task := os.Args[1]
	fmt.Printf("map reduce task: %s\n", task)

	masterInit(task)
}