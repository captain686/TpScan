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

var (
	url      string
	filePath string
)

func init() {
	flag.StringVar(&url, "u", "", "Target Url")
	flag.StringVar(&filePath, "f", "", "TargetFile Path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of question:\n")
		flag.PrintDefaults()
	}
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
	flag.Parse()
	if flag.NFlag() == 1 {
		if url != "" {
			util.ExpExecutor(url)
		}
		if filePath != "" {
			for target := range fileRead(filePath) {
				util.ExpExecutor(*target)
			}
		}
	}
	if flag.NFlag() > 1 {
		fmt.Fprintf(os.Stderr, "Usage of question:\n")
		flag.PrintDefaults()
	}
}
