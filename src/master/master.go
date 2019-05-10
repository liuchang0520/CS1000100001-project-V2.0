package main

import (
	"fmt"
	"os"
	// "google.golang.org/grpc"
	// "net"
	"strconv"
	"io/ioutil"
	// "strings"
	"bufio"
	"time"
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
	inputDir string
}

type WorkArgs struct {
	Task string
	Stage string //map or reduce

	//for map, it is an input file; for reducer, it is a directory containing the intermediate files
	InputFile string
}
type WorkerRes struct {

}
type RegisterArgs struct {
	Port string
}
type MasterRes struct {
	RCnt int // number of reducers we want to use
	InputDir string
}

func (master *Master) RegisterWorker(args *RegisterArgs, res *MasterRes) error {
	fmt.Println(args.Port, " registering")

	//create client in master for this worker
	fmt.Printf("dial worker at port: %s \n", args.Port)
	masterClient, err := rpc.DialHTTP("tcp", "localhost:" + args.Port)
	if err != nil {
		fmt.Printf("failed to dial worker at port: %s \n", args.Port)
		defer log.Fatal(err)
		return err
	}
	fmt.Println("master client created")
	master.client[args.Port] = masterClient

	master.workerChan <- args.Port
	fmt.Printf("worker at port: %s registered\n", args.Port)

	res.RCnt = master.rCnt
	res.InputDir = master.inputDir
	return nil
}

func masterInit(task, inputDir, rCnt string) *Master {
	master := new(Master)
	master.workerChan = make(chan string, c.MAX_WORKER)
	master.task = task
	cnt, err := strconv.Atoi(rCnt)
	if err != nil {
		log.Fatal("parse number of reducer failed", err)
	}
	master.rCnt = cnt
	master.inputDir = inputDir
	master.input = c.GetInputF(inputDir)
	master.client = make(map[string]*rpc.Client)
	master.closeChan = make(chan bool)
	// master.workerCloseChan = make(chan bool)

	// l, err := net.Listen("tcp", c.MASTER_PORT)
	// if err != nil {
	// 	log.Fatal("master failed to listen: ", err)
	// }
	// master.listener = l

	fmt.Println("master rpc register")
	if err = rpc.Register(master); err != nil {
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

func runTask(master *Master) {
	fmt.Println("start running tasks...")
	go func() {
		//this is just for Project presentation
		fmt.Println("wait 10 seconds to get task started")
		for i := 1; i <= 10; i++ {
			time.Sleep(time.Second)	
			fmt.Println(10 - i + 1)
		}

		assignTask(master, c.MAP)
		assignTask(master, c.REDUCE)
		aggregate() //aggregate reducer output files into a single file
		finish(master)
	}()
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
		if err := c.CreateInterDir(master.rCnt); err != nil {
			log.Fatal(err)
		}
		taskCnt = len(master.input)
	} else if stage == c.REDUCE {
		if err := c.CreateOutputDir(); err != nil {
			log.Fatal(err)
		}
		taskCnt = master.rCnt
	} else {
		log.Fatal("unknown stage: ", stage)
	}

	var wg sync.WaitGroup //use this to wait all tasks to be finished

	for i := 0; i < taskCnt; i++ { //assign task
		wg.Add(1)
		fmt.Printf("start assigning %s task %d\n", stage, i)
		go func (taskIndex int) {
			temp := 0
			for temp < c.MAX_TEMP {
				//get an idle worker
				w := <- master.workerChan
				var err error

				var inF string
				if stage == c.MAP {
					inF = master.input[taskIndex]
				} else {
					inF = fmt.Sprintf("%d", taskIndex)
				}
				fmt.Printf("call worker at port: %s \n", w)
				if err = master.client[w].Call("Worker.Work", &WorkArgs{Stage: stage, Task: master.task, InputFile: inF}, &WorkerRes{}); err == nil {
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
			log.Fatal(fmt.Sprintf("in stage %s, %d task cannot be finished\n", stage, taskIndex))
		} (i)
	}

	wg.Wait() //use this to wait all tasks to be finished
	fmt.Printf("%s stage complete\n", stage)
}

func aggregate() {
	//read output files from reducers
	files, err := ioutil.ReadDir(fmt.Sprintf("../%s/", c.OUTPUT_DIR))
	if err != nil {
		fmt.Println("aggregation output failed: ", err)
		return 
	}

	//create aggregated output file
	outF, err := os.Create(fmt.Sprintf("../%s", c.AGGREGATED_OUT_F))
	if err != nil {
		fmt.Println("aggregation output failed: ", err)
		return
	}
	defer outF.Close()

	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("../%s/", c.OUTPUT_DIR) + file.Name())
		if err != nil {
			fmt.Println("aggregation output failed: ", err)
			return
		}
    	defer f.Close()

	    scanner := bufio.NewScanner(f)
	    for scanner.Scan() {
	    	if _, err = outF.WriteString(scanner.Text() + "\n"); err != nil {
				fmt.Println("aggregation output failed: ", err)
				return
			}
	    }

	    if err = scanner.Err(); err != nil {
	    	fmt.Println("aggregation output failed: ", err)
	        return
	    }
	}

	fmt.Printf("output files all have been aggregated to %s\n", c.AGGREGATED_OUT_F)
}

func finish(master *Master) {
	for _, ct := range(master.client) {
		if err := ct.Call("Worker.Close", &WorkArgs{}, &WorkerRes{}); err != nil {
			log.Fatal(err)
		}
	}
	master.closeChan <- true
}

func main() {
	//args[1] indicates the map reduce task
	//args[2] indicates the input dir
	//args[3] indicates the number of reducers
	if len(os.Args) != 4 {
		log.Fatal("pls specify the mapreduce task, input file directory, and number of reducers", " args: ", os.Args)
	}
	task := os.Args[1]
	fmt.Printf("map reduce task: %s\n", task)
	if _, ok := c.FuncMap[task]; !ok {
		log.Fatal("invalid task, pls check the task name", fmt.Sprintf("valid task names are: %s\n", c.GetValidTask()))
	}

	master := masterInit(task, os.Args[2], os.Args[3])
	runTask(master)
	<- master.closeChan
	// master.listener.Close()
	//send shut down signal to workers
	// <- master.workerCloseChan
	fmt.Println("master work complete")
}