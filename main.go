package main

import (
	"flag"
	"fmt"
	"github.com/yngwiewang/ansible-go/internal"
	"sync"
)

var (
	hostsFile  string
	module     string
	cmd        string
	//outputFile string
	fileSrc    string
	fileDest   string
	fileMask   string
	hosts      []Host
	hostsCount int
	resultChan chan string
	wg         sync.WaitGroup
	err        error
)

func init() {
	hostsFile = "hosts.csv"
	flag.StringVar(&hostsFile, "i", hostsFile, "csv file including hosts information")
	flag.StringVar(&module, "m", "shell", "module, could be shell or copy, default shell")
	flag.StringVar(&cmd, "c", "", "bash command to execute")
	flag.StringVar(&fileSrc, "src", "", "source of the file to copy")
	flag.StringVar(&fileDest, "dest", "", "destination of the file to copy")
	flag.StringVar(&fileMask, "mask", "0744", "mask of the file to copy, default 0744")

	hosts, err = ReadCsv(hostsFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hostsCount = len(hosts)
	resultChan = make(chan string, hostsCount)

	wg = sync.WaitGroup{}
	wg.Add(hostsCount)
}

func main() {

	flag.Parse()
	//if module == "" {
	//	module = "shell"
	//}

	modules := []string{"shell", "copy"}
	if internal.Contain(modules, module) == false {
		fmt.Println("You must specify a module")
	}
	switch module {
	case "shell":
		if cmd == "" {
			fmt.Println("You must specify a bash command")
			return
		}

		for _, h := range hosts {

			go func(h Host, cmd string) {
				defer wg.Done()
				resultChan <- h.Exec(cmd)
			}(h, cmd)
		}

		for i := 0; i < hostsCount; i++ {
			fmt.Println(<-resultChan)
		}
		wg.Wait()
	case "copy":
		if fileSrc == "" {
			fmt.Println("You must specify a source file to copy")
			return
		}
		if fileDest == "" {
			fmt.Println("You must specify a destination to copy")
			return
		}
		for _, h := range hosts {

			go func(h Host, src, dest, mask string) {
				defer wg.Done()
				resultChan <- h.Copy(src, dest, mask)
			}(h, fileSrc, fileDest, fileMask)
		}

		for i := 0; i < hostsCount; i++ {
			fmt.Println(<-resultChan)
		}
		wg.Wait()

	default:
		fmt.Println("You must specify a module")
	}

}
