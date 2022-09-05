package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/captain686/TpScan/util"
	"os"
)

func init() {
	banner := `
	████████╗██████╗ ███████╗ ██████╗ █████╗ ███╗   ██╗
	╚══██╔══╝██╔══██╗██╔════╝██╔════╝██╔══██╗████╗  ██║
	   ██║   ██████╔╝███████╗██║     ███████║██╔██╗ ██║
	   ██║   ██╔═══╝ ╚════██║██║     ██╔══██║██║╚██╗██║
	   ██║   ██║     ███████║╚██████╗██║  ██║██║ ╚████║
	   ╚═╝   ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝`
	fmt.Println(banner)
}

func getArgs() (*string, string) {
	var (
		url      string
		filePath string
	)
	flag.StringVar(&url, "u", "", "Target Url")
	flag.StringVar(&filePath, "f", "", "TargetFile Path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of question:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NFlag() == 1 {
		if url != "" {
			//fmt.Println(url)
			return &url, "url"
		}
		if filePath != "" {
			fmt.Println(filePath)
			return &filePath, "filePath"
		}
	}
	if flag.NFlag() > 1 || flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of question:\n")
		flag.PrintDefaults()
	}
	return nil, ""
}

func fileRead(path string) <-chan *string {
	ch := make(chan *string)
	fp, err := os.Open(path)
	if util.IfErr(err) {
		fmt.Println(err) //打开文件错误
		return nil
	}
	buf := bufio.NewScanner(fp)
	go func() {
		for {
			if !buf.Scan() {
				break //文件读完了,退出for
			}
			line := buf.Text() //获取每一行
			ch <- &line
		}
		close(ch)
	}()
	return ch
}

func main() {
	value, valueType := getArgs()
	if valueType == "url" {
		util.ExpExecutor(*value)
	}
	if valueType == "filePath" {
		for target := range fileRead(*value) {
			util.ExpExecutor(*target)
		}
	}
}
