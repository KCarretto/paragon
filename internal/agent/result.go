package agent

// Result stores and displays the result of running a task.
type Result struct {
	TaskID    string
	TimeStart int
	TimeEnd   int
	Stdout    string
	Stderr    string
	ExitCode  int
	Err       bool
}
