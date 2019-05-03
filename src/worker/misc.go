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
	return fmt.Sprintf("../%s-%d/", c.INTERMEDIATE_DIR, i)
}

func getOutputF(index string) string {
	return fmt.Sprintf("../%s/%s-%s", c.OUTPUT_DIR , c.OUTPUT_F_PREFIX, index)
}

func createInterDir(rCnt int) error {
	for i := 0; i < rCnt; i++ {
		if err = os.Mkdir(getInteDir(i), 0777); err != nil {
			return err
		}
	}

	return nil
}

func createOutputDir() error {
	if err := os.Mkdir(fmt.Sprintf("../%s/", c.OUTPUT_DIR), 0777); err != nil {
		return err
	}

	return nil
}