# Overview

This project simulates a map reduce system:

Different processes simulate master and workers. Task assignment could be achieved by RPC between master and worker processes.

Implemented very basic fault tolerance mechanism by Go Channel: https://gobyexample.com/channels , in which idle worker process ports are stored.

Other features:  
Worker registration: by RPC

More detailed idea/implementation is in proposal.txt.

# Adding your own map reduce task:

So far I have only implemented a word count task. To add new map reduce task, pls refer to FuncMap, wordCountMapFunc, wordCountReduceFunc in common.go. Don't forget to include the task name when you start master process.

# Steps to run this project in Shell:

* golang environment setup: https://golang.org/doc/install

* set GOPATH to the root directory of this project: cd to root dir , then do: export "GOPATH=$PWD". To check whether GOPATH has changed, do: go env

* create an input file directory under the src dir, where you can put input files. I have created a wordcount_input folder for reference. 
You may have to change the permission of this input folder.

* open separate shell windows for master and workers setup:

* * For master process, cd to src/master, do:  go run master.go [task_name] [input_dir_name] [number of expected reducers] 
* * * e.g. go run master.go wordCount wordcount_input 4 ==> This will run a word count task using 4 reducers at port 5450.
 You could change the master process port number in MASTER_PORT in common.go

* * For worker processes, cd to src/worker, do: go run worker.go [any_available_port_number_other_than_5450]
* * * e.g. go run worker.go 5451

* After one map reduce work, pls remove the generated files/dir before doing a new one.