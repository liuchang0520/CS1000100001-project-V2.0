package commmon

const (
	MAP = "map"
	REDUCE = "reduce"

	MASTER_PORT = "5450"
	MAX_WORKER = 10

	MAX_TEMP = 10 //maximum number of task temp before announcing failure

	INTERMEDIATE_DIR := "inte"
	INTERMEDIATE_F_PREFIX := "inter-"

	OUTPUT_F_PREFIX := "out"
	OUTPUT_DIR := "outdir"
)


type KV struct {
	k string
	v string
}

type WorkArgs struct {
	Task string
	Stage string //map or reduce

	//for map, it is an input file; for reducer, it is a directory containing the intermediate files
	InputFile string 
}
type WorkerRes struct {

}
type RegisterArgs struct {
	Port string
}
type MasterRes struct {
	RCnt int // number of reducers we want to use
}

//TODO: map[functionName]function e.g. "wordCount": func wordCount()....