# Overview

This project simulates a map reduce system:

Different processes simulate master and workers. Task assignment could be achieved by RPC between master and worker processes.

Implemented very basic fault tolerance mechanism by Go Channel: https://gobyexample.com/channels , in which idle worker process ports are stored.

Other features:  
worker registration: by RPC

More detailed idea/implementation is in proposal.txt.



# Steps to run this project in Shell:

* golang environment setup: https://golang.org/doc/install

* set GOPATH to the root directory of this project: cd to root dir , then do: export "GOPATH=$PWD". To check whether GOPATH has changed, do: go env

* create a input directory in the root dir, where you can put arbitray files.

* open separate shell windows for master and workers,

* * dsds
* * dasda



go run master.go wordCount wordcount_input 4



go run worker.go 5452 


* clean the generated file

