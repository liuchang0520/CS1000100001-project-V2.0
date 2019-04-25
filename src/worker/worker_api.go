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
	0) read corresponding input file
	1) transform the file content to k:v pairs
		input: 
		k: input file name
		v: content

	2) for each k, do hash(k) % r, to get which intermediate file dir to put to
	   put all intermediate files for a reducer into a folder
	*/
	
	/*
	reduce:

	1 output file per reducer
	0) read the ith directory, in which intermediate files for this reducer is stored
	1) transform to k: v0,v1,v2....combiner
	2) call reduce function for each k, v
		input: 
		k: the key from intermediate files
		v: the aggregated value list
	3) for each output of reduce function, put it into the output file
	*/

	


	return nil
}

//shutdown worker
func (worker *Worker) Close(args *c.WorkArgs, res *c.WorkerRes) error {
	worker.shutDownChan <- true
	fmt.Printf("shutdown signal received for worker: %v\n", worker.port)
	return nil
}
