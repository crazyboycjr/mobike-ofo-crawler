package main

import (
	"os"
	"runtime"
	"flag"

	"github.com/crazyboycjr/mobike-ofo-crawler/mobike"
	"github.com/crazyboycjr/mobike-ofo-crawler/ofo"
)

func main() {
	runtime.GOMAXPROCS(1) // Do not change

	var token, ofile string
	flag.StringVar(&token, "-token", "token.txt", "Specify the file saving your tokens")
	flag.StringVar(&ofile, "-output", "data.txt", "Specify the output file")

	flag.Parse()

	fout, err := os.OpenFile(ofile, os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	wait := make(chan int)
	go mobike.Run(token, fout)
	go ofo.Run(fout)
	<-wait
}
