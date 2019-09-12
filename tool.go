package tool

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	read(stdout, true, &success)
	read(stderr, false, &success)
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

//DeCompress 解压zip文件
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + "/" + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
