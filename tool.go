package tool

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

//Execute 执行命令
func Execute(args []string) bool {
	fmt.Println("当前执行命令: ", args)
	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return false
	}
	success := false
	go read(stdout, true, &success)
	go read(stderr, false, &success)
	cmd.Wait()
	return success
}

func read(rc io.ReadCloser, isOut bool, success *bool) {

	reader := bufio.NewReader(rc)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if !isOut {
			*success = false
		} else {
			*success = true
		}
		fmt.Print(line)
	}

}
