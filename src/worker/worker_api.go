package main

import (
	"fmt"
	"io/ioutil"
	"os"
	// "log"
	"errors"
	// "net"
	// "net/rpc"
	// "net/http"
	c "common"
)

func mapTask(task, input string, rCnt int, mapFunc func(string, string) c.KV) error {
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

	var content []byte
	if content, err := ioutil.ReadFile(input); err != nil {
		return err
	}

	words := splitWords(string(content))

	//buffer the intermediate files
	//rCnt intermediate files per mapper
	inteFile := make([]string, rCnt)
	for _, word := range(words) { 
		//get a key value pair out of KV defined in common.go 
		kv := mapFunc(input, word)

		//using hash to determine the intermediate file
		inteIndex := getHashCode(kv.k) % rCnt
		inteFile[inteIndex] += kv.k + ":" + kv.v + "\n"
	}

	//write buffered data to intermediate file
	for i, str := inteFile {
		inteName := getInteDir(i) + c.INTERMEDIATE_F_PREFIX + input
		f, err := os.Create(inteName)
		if err != nil {
			return err
		}

		if _, err = f.WriteString(str[:len(str) - 1]); err != nil {
			return err
		}

		f.Close()
	}

	fmt.Printf("map task for %s in task %s is complete\n", input, task)
}

func reduceTask(task, input string, rCnt int, reduceFunc func(string, string) c.KV) error {
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

	//create output file
	outputFile := getOutputF(input)
	outF, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outF.Close()

	//TODO: 0) read ......
}

//implement worker's rpc api
//this rpc api is used to call worker to do a specific map/reduce task
func (worker *Worker) Work(args *c.WorkArgs, res *c.WorkerRes) error {
	fmt.Printf("worker: %s doing %s of task - %s\n", worker.port, args.Stage, args.Task)

	//get the corresponding map and reduce functions for this task
	mrFunc := funcMap[args.Task]

	if args.Stage == c.MAP {
		return mapTask(args.Task, args.InputFile, worker.rCnt, mrFunc[0])
	}

	if args.Stage == c.REDUCE {
		return reduceTask(args.Task, args.InputFile, worker.rCnt, mrFunc[1])
	}
	
	return errors.New(fmt.Sprint("invalid stage: %s\n", args.Stage))
}


//shutdown worker
func (worker *Worker) Close(args *c.WorkArgs, res *c.WorkerRes) error {
	worker.shutDownChan <- true
	fmt.Printf("shutdown signal received for worker: %v\n", worker.port)
	return nil
}
