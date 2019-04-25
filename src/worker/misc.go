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

func splitWords(str string) []string {

}

func getHashCode(key string) int {

}


func getInteDir(i int) string {
	return fmt.Sprintf("/%s-%d/", c.INTERMEDIATE_DIR, i)
}

func getOutputF(index string) string {
	return fmt.Sprintf("/%s/%s-%s", c.OUTPUT_DIR , c.OUTPUT_F_PREFIX, index)
}