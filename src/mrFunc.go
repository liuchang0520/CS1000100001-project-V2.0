package main

import (
	"fmt"
	"io/ioutil"
	// "os"
	// "log"
	"errors"
	// "net"
	// "net/rpc"
	// "net/http"
	c "common"
)

var (
	funcMap := map[string][]func(string, string) c.KV
    	{"wordCount": [2]func(string, string) c.KV {wordCountMapFunc, wordCountReduceFunc}} 
)

//for wordCount
func wordCountMapFunc(key, val string) c.KV {

}

func wordCountReduceFunc(key, val string) c.KV {

}