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

type MRFunc struct {
	MF func(string, string) c.KV
	RF func(string, []string) c.KV
}

var (
	funcMap := map[string]MRFunc
    	{"wordCount": MRFunc{wordCountMapFunc, wordCountReduceFunc}} 
)

//for wordCount
func wordCountMapFunc(key, val string) c.KV {

}

func wordCountReduceFunc(key, val string) c.KV {

}