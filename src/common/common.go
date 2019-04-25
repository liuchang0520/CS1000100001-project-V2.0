package commmon

const (
	MAP = "map"
	REDUCE = "reduce"

	MASTER_PORT = "5450"
	MAX_WORKER = 10

	MAX_TEMP = 10 //maximum number of task temp before announcing failure

	outputDir := ""
)

type WorkArgs struct {
	Task string
	Stage string //map or reduce

	//for map, it is an input file; for reducer, it is a directory containing the intermediate files
	inputFile string 


}
type WorkerRes struct {

}
type RegisterArgs struct {
	Port string
}
type MasterRes struct {
	
}

//TODO: map[functionName]function e.g. "wordCount": func wordCount()....