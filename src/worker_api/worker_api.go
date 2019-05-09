package worker_api

import (
	"fmt"
	"bufio"
	"io/ioutil"
	"strings"
	"os"
	"strconv"
	"log"
	// "errors"
	// "net"
	// "net/rpc"
	// "net/http"
	c "common"
)

func MapTask(task, inputDir, input string, rCnt int, mapFunc func(string, string) c.KV) error {
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
	var err error
	if content, err = ioutil.ReadFile(fmt.Sprintf("../%s/", inputDir) + input); err != nil {
		fmt.Println(err)
		return err
	}

	words := c.SplitWords(string(content))

	//buffer the intermediate files
	//rCnt intermediate files per mapper
	inteFile := make([]string, rCnt)
	for _, word := range(words) { 
		//get a key value pair out of KV defined in common.go 
		kv := mapFunc(input, word)

		//using hash to determine the intermediate file
		inteIndex := c.GetHashCode(kv.K) % rCnt
		inteFile[inteIndex] += kv.K + c.KV_SEP + kv.V + "\n"
	}

	//write buffered data to intermediate file
	for i, str := range inteFile {
		if len(str) == 0 {
			continue
		}
		inteName := c.GetInteDir(i) + c.INTERMEDIATE_F_PREFIX + input
		f, err := os.Create(inteName)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if _, err = f.WriteString(str[:len(str) - 1]); err != nil {
			fmt.Println(err)
			return err
		}

		f.Close()
	}

	fmt.Printf("map task for %s in task %s is complete\n", input, task)
	return nil
}

func ReduceTask(task, input string, rCnt int, reduceFunc func(string, []string) c.KV) error {
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
	outputFile := c.GetOutputF(input)
	outF, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer outF.Close()

	//build a dictionary to group key-value pair from intermediate files based on keys
	kvMap := make(map[string][]string)

	//read intermediate files directory for this reducer
	i, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal("parse intermediate file failed: ", input)
	}
	interDir := c.GetInteDir(i);
	files, err := ioutil.ReadDir(interDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		f, err := os.Open(interDir + file.Name())
		if err != nil {
			fmt.Println(err)
			return err
		}
    	defer f.Close()

	    scanner := bufio.NewScanner(f)
	    for scanner.Scan() {
	    	kv := strings.Split(scanner.Text(), c.KV_SEP)
	    	if len(kv) != 2 {
	    		continue
	    		// log.Fatal("cannot parse intermediate file: ", file.Name())
	    	}
	    	key, val := kv[0], kv[1]
	    	vList := []string{}
	        if tempList, ok := kvMap[key]; ok {
	        	vList = tempList
	        }
	        kvMap[key] = append(vList, val)
	    }

	    if err = scanner.Err(); err != nil {
	    	fmt.Println(err)
	        return err
	    }
	    fmt.Printf("intermediate file: %s has been processed\n", file.Name())
	}

	//call reduce function per key
	for k, v := range kvMap {
		kvRes := reduceFunc(k, v)
		if _, err = outF.WriteString(kvRes.K + ":" +kvRes.V + "\n"); err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Printf("reduce task for %s in task %s is complete\n", input, task)
	return nil
}