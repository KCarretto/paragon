package agent_test

//go:generate mockgen -destination=io.gen_test.go -package=agent_test io Writer,WriteCloser
//go:generate mockgen -destination=agent.gen_test.go -package=agent_test github.com/kcarretto/paragon/agent Sender,Receiver
