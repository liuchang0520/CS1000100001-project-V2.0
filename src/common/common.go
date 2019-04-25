package commmon

const (
	MAP = "map"
	REDUCE = "reduce"

	MASTER_PORT = "5450"
	MAX_WORKER = 10

	MAX_TEMP = 10 //maximum number of task temp before announcing failure
)

type WorkArgs struct {
	Task string
}
type WorkerRes struct {

}
type RegisterArgs struct {
	Port string
}
type MasterRes struct {
	
}