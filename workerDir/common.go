package main

const (
	MASTER_PORT = "5450"
	MAX_WORKER = 10
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
