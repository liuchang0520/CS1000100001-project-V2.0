package main

import (
	"fmt"
	"os"
	"io/ioutil"
	// "google.golang.org/grpc"
	// "net"
	"log"
	"net/http"
	c "common"
	"net/rpc"
)

func getInputF(inputDir string) []string {
	files := []string{}

	fs, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal("read input dir failed: ", err)
	}
	for _, f := range fs {
		files = append(files, f.Name())

	return files
}

func getValidTask() string {
	res := ""

	for task, _ := range(funcMap) {
		res += task + "\n"
	}

	return res
}


