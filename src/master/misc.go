package main

import (
	"fmt"
	"os"
	// "google.golang.org/grpc"
	// "net"
	"log"
	"net/http"
	c "common"
	"net/rpc"
)

func getInputF(input string) []string {
	
}

func getValidTask() string {
	
}


func getInteDir(i int) string {
	return fmt.Sprintf("/%s-%d/", c.INTERMEDIATE_DIR, i)
}