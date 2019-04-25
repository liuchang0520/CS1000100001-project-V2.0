package main

import (
	"fmt"
	// "os"
	// "google.golang.org/grpc"
	// "net"
	// "log"
	// "net/http"
	c "common"
	"net/rpc"
)

//implement master's rpc api
//this rpc api is used to register workers
func (master *Master) RegisterWorker(args *c.RegisterArgs, res *c.MasterRes) error {
	fmt.Println(args.Port, " registering")

	//create client in master for this worker
	masterClient, err := rpc.DialHTTP("tcp", "localhost:" + args.Port)
	if err != nil {
		defer log.Fatal(err)
		return err
	}
	fmt.Println("master client created")
	master.client[args.Port] = masterClient

	master.workerChan <- args.Port
	fmt.Printf("worker at port: %s registered\n", args.Port)
	return nil
}