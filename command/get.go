package command

import "fmt"

type getCmd struct {
	err     error
	succeed bool
	result  string
}

func (e *getCmd) Exec(cmd *Command) {
	fmt.Println(cmd)
	e.succeed = true
	e.result = "Command,get, exec successfully!"
}

func (e *getCmd) Succeed() bool {
	return e.succeed
}

func (e *getCmd) Result() string {
	return e.result
}

func (e *getCmd) Error() error {
	return e.err
}

func (e *getCmd) flush() {
	e.result = ""
	e.err = nil
	e.succeed = false
}
