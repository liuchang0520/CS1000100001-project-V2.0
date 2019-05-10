package commmon

import (
	"fmt"
	"os"
	"io/ioutil"
	// "errors"
	// "google.golang.org/grpc"
	// "net"
	"log"
	// "net/http"
	// "net/rpc"
	// "strings"
)

//map reduce function
type MRFunc struct {
	MF func(string, string) KV
	RF func(string, []string) KV
}

type KV struct {
	K string
	V string
}

const (
	MAP = "map"
	REDUCE = "reduce"

	MASTER_PORT = "5450"
	MAX_WORKER = 10

	MAX_TEMP = 10 //maximum number of task temp before announcing failure

	INTERMEDIATE_DIR = "inte"
	INTERMEDIATE_F_PREFIX = "inter-"

	OUTPUT_F_PREFIX = "out"
	OUTPUT_DIR = "outdir"

	AGGREGATED_OUT_F = "output.txt"

	KV_SEP = ":"
	// TEXT_SEP = " "

	ONE = "1"
)

//map which stores valid map reduce task functions
var (
	FuncMap = map[string]MRFunc{"wordCount": MRFunc{wordCountMapFunc, wordCountReduceFunc}} 
)

//following is the map reduce functions
//for wordCount
func wordCountMapFunc(key, val string) KV {
	return KV{K: val, V: ONE}
}

func wordCountReduceFunc(key string, val []string) KV {
	return KV{K: key, V: fmt.Sprintf("%d", len(val))}
}



//misc functions

func GetInputF(inputDir string) []string {
	files := []string{}

	fs, err := ioutil.ReadDir(fmt.Sprintf("../%s/", inputDir))
	if err != nil {
		log.Fatal("read input dir failed: ", err)
	}
	for _, f := range fs {
		files = append(files, f.Name())
	}
	return files
}

func GetValidTask() string {
	res := ""

	for task, _ := range(FuncMap) {
		res += task + "\n"
	}

	return res
}

// func SplitWords(str string) []string {
// 	return strings.Split(str, TEXT_SEP)
// }

//get the hash code like Java String.hashCode()
func GetHashCode(key string) int {
	ans := 0
	base := 1

	for i := len(key) - 1; i >= 0; i-- {
		ans += base * int(key[i])
		base *= 31
	}

	return ans
}

func GetInteDir(i int) string {
	return fmt.Sprintf("../%s-%d/", INTERMEDIATE_DIR, i)
}

func GetOutputF(index string) string {
	return fmt.Sprintf("../%s/%s-%s", OUTPUT_DIR , OUTPUT_F_PREFIX, index)
}

func CreateInterDir(rCnt int) error {
	for i := 0; i < rCnt; i++ {
		if err := os.Mkdir(GetInteDir(i), 0777); err != nil {
			return err
		}
	}

	return nil
}

func CreateOutputDir() error {
	if err := os.Mkdir(fmt.Sprintf("../%s/", OUTPUT_DIR), 0777); err != nil {
		return err
	}

	return nil
}