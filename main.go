package main

import (
	"os"
	"runtime"
	"flag"
	"strings"
	"fmt"

	"github.com/crazyboycjr/mobike-ofo-crawler/mobike"
	"github.com/crazyboycjr/mobike-ofo-crawler/ofo"
)

type moduleList []string

func (m *moduleList) String() string {
	return strings.Join(*m, ",")
}

func (m *moduleList) Set(value string) error {
	for _, module := range strings.Split(value, ",") {
		*m = append(*m, module)
	}
	return nil
}

func contains(module []string, target string) bool {
	for _, s := range module {
		if s == target {
			return true
		}
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(1) // Do not change

	var token, ofile string
	flag.StringVar(&token, "token", "token.txt", "Specify the file saving your tokens")
	flag.StringVar(&ofile, "output", "data.txt", "Specify the output file")

	var module moduleList
	flag.Var(&module, "module", "Comma-separated list of modules to use")

	var conMobike, conOfo int
	flag.IntVar(&conMobike, "Cmobike", 10, "the concurrency number of mobike module")
	flag.IntVar(&conOfo, "Cofo", 1, "the concurrency number of ofo module")

	flag.Parse()

	fout, err := os.OpenFile(ofile, os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println("you specify the modules:", module)

	if len(module) == 0 {
		module = []string{"mobike", "ofo"}
		fmt.Println("all modules, will be used")
	}

	wait := make(chan int)
	if contains(module, "mobike") {
		go mobike.Run(token, fout, conMobike)
	}
	if contains(module, "ofo") {
		go ofo.Run(fout, conOfo)
	}
	<-wait
}
