package agent

// Tasks can be run, which return a result
type Task interface {
	Run() Result
}
