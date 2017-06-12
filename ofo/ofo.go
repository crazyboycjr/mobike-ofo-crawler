package ofo

import (
	"os"
	"log"
	"time"

	"net/http"
	"encoding/json"
	"math/rand"

	"github.com/crazyboycjr/mobike-ofo-crawler/utility"
)

const ofoUrl string = "https://open.ofo.so/v1/near/bicycle"

func setReqHeader(req *http.Request) {
	req.Header.Set("Host", "open.ofo.so")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", "OneTravel/5.0.16 (iPhone; iOS 10.2; Scale/3.00)")
	req.Header.Set("Accept-Language", "en-CN;q=1, zh-Hans-CN;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
}

type Position struct {
	lng, lat float64
}

type Response struct {
	pos Position
	text []byte
}

func consume(chreq chan Position, chres chan Response) {

	for pos := range chreq {
		timeout := time.Duration(1 * time.Second)
		client := &http.Client{
			Timeout: timeout,
		}
		lng, lat := pos.lng, pos.lat

		form := make(map[string] interface{})
		form["appKey"] = "didi"
		form["mapType"] = "1"
		form["longitude"] = lng
		form["latitude"] = lat
		form["datatype"] = "101"
		form["appversion"] = "5.0.16"
		data, err := json.Marshal(form)
		if err != nil {
			panic(err)
		}
		//log.Println(string(data))

		text := utility.ReadAll(client, ofoUrl, string(data), setReqHeader)

		log.Println(string(text))
		chres <- Response{pos, text}
	}
}

var dic = make(map[string] int)
var fout *os.File

func produce(chreq chan Position, chres chan Response) {

	for res := range chres {
		text := res.text
		x := res.pos.lng
		y := res.pos.lat

		var j map[string] interface{}
		json.Unmarshal(text, &j)
		if j["body"] != nil {
			for _, o := range j["body"].(map[string] interface{})["bicycles"].([]interface{}) {
				ob := o.(map[string] interface{})
				id := ob["bicycleNo"].(string)
				if _, ok := dic[id]; !ok || (utility.NowMillisec() - dic[id]) > 5 * 60 * 1000 {
					dic[id] = utility.NowMillisec()
					//write file
					ob["ts"] = dic[id]
					out, _ := json.Marshal(ob)
					fout.Write(out)
					fout.WriteString("\n")
				}
				//xx := ob["longitude"].(float64)
				//yy := ob["latitude"].(float64)
				//log.Println(xx, yy)
			}
		}
		log.Println("ofo handle", x, y)
	}
}

func Run(outfile *os.File, conNum int) {
	rand.Seed(time.Now().UnixNano())
	fout = outfile

	const R float64 = 1000.0
	const K float64 = R * 0.001

	pos := []Position{}
	for x := 120.85000000; x < 122.2; x += 0.008984 * K {
		for y := 30.6666666; y < 31.88333333; y += 0.008984 * K {
			pos = append(pos, Position{x, y})
		}
	}
	log.Println("The number of positions = ", len(pos))
	tot := len(pos)
	for i := 1; i < tot; i++ {
		j := rand.Intn(i)
		pos[i], pos[j] = pos[j], pos[i]
	}

	chreq := make(chan Position)
	chres := make(chan Response)

	for i := 0; i < conNum; i++ {
		go consume(chreq, chres)
		go produce(chreq, chres)
	}
	for {
		for _, p := range pos {
			chreq <- Position{p.lng, p.lat}
		}
		time.Sleep(time.Second * 10)
		log.Println("Round finished")
		break
	}
}

