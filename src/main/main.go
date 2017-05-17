package main

import (
	/*
	"log"
	"net/http"
	"io/ioutil"
	"time"
	_"bytes"
	"strings"
	"net/url"
	"encoding/json"
	"os"
	_"fmt"
	"strconv"
	"math/rand"
	*/
	"os"
	"runtime"
	"mobike"
	"ofo"
	"flag"
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
