package main

import (
	"fmt"
	// "os"
	// "log"
	// "net"
	// "net/rpc"
	// "net/http"
	c "common"
)

//implement worker's rpc api
//this rpc api is used to call worker to do a specific map/reduce task
func (worker *Worker) Work(args *c.WorkArgs, res *c.WorkerRes) error {
	/*
	map:
	put all intermediate files for a reducer into a folder
	*/
	
	/*
	*/
	


	return nil
}

//shutdown worker
func (worker *Worker) Close(args *c.WorkArgs, res *c.WorkerRes) error {
	worker.shutDownChan <- true
	fmt.Printf("shutdown signal received for worker: %v\n", worker.port)
	return nil
}
