package mobike

import (
	"os"
	"log"
	"time"
	"strconv"

	"net/http"
	"net/url"
	"encoding/json"
	"math/rand"

	"github.com/crazyboycjr/mobike-ofo-crawler/utility"
)

const conNum = 10
const mobikeUrl string = "https://mwx.mobike.com/mobike-api/rent/nearbyBikesInfo.do"

var tokens []Token

func setReqHeader(req *http.Request) {
	no := rand.Intn(len(tokens))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "MicroMessenger/6.5.7.1041 NetType/WIFI Language/en")
	req.Header.Set("mobileno", tokens[no].Mobileno)
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("time", string(utility.NowMillisec()))
	req.Header.Set("Aceept", "*/*")
	req.Header.Set("open_src", "list")
	req.Header.Set("Referer", "https://servicewechat.com/wx80f809371ae33eda/23/")
	req.Header.Set("platform", "4")
	req.Header.Set("citycode", "021")
	req.Header.Set("lang", "zh")
	req.Header.Set("accesstoken", tokens[no].Accesstoken)
	req.Header.Set("eption", "")
	req.Header.Set("charset", "utf-8")
}

type Position struct {
	lng, lat, radius float64
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

		form := url.Values{}
		form.Set("errMsg", "getLocation:ok")
		form.Set("longitude", strconv.FormatFloat(lng, 'f', -1, 64))
		form.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
		form.Set("accuracy", string(30))
		form.Set("citycode", "021")
		form.Set("speed", "0")
		data := form.Encode()

		text := utility.ReadAll(client, mobikeUrl, data, setReqHeader)

		log.Println(string(text))
		chres <- Response{pos, text}
	}
}

var dic = make(map[string] int)
var fout *os.File

func produce(chreq chan Position, chres chan Response) {

	for res := range chres {
		text := res.text
		radius := res.pos.radius
		x := res.pos.lng
		y := res.pos.lat

		k := radius * 0.001
		delta := k * 0.008984
		x1 := x - delta
		y1 := y - delta
		x2 := x + delta
		y2 := y + delta
		inRangeCount := 0

		dist := 0
		var j map[string] interface{}
		json.Unmarshal(text, &j)
		if j["object"] != nil {
			for _, o := range j["object"].([]interface{}) {
				ob := o.(map[string] interface{})
				id := ob["bikeIds"].(string)
				if _, ok := dic[id]; !ok || (utility.NowMillisec() - dic[id]) > 5 * 60 * 1000 {
					dic[id] = utility.NowMillisec()
					//write file
					ob["ts"] = dic[id]
					out, _ := json.Marshal(ob)
					fout.Write(out)
					fout.WriteString("\n")
				}
				xx := ob["distX"].(float64)
				yy := ob["distY"].(float64)
				d, err := strconv.Atoi(ob["distance"].(string))
				if err == nil && dist < d {
					dist = d
				}
				//log.Println(xx, yy)
				if x1 < xx && xx < x2 && y1 < yy && yy < y2 {
					inRangeCount++
				}
			}
		}
		log.Println("handle", x, y, inRangeCount, radius, dist)

		if inRangeCount >= 23 || dist <= 100 {
			go func () {
				chreq <- Position{(x1 + x) * 0.5, (y1 + y) * 0.5, radius * 0.5}
				chreq <- Position{(x1 + x) * 0.5, (y + y2) * 0.5, radius * 0.5}
				chreq <- Position{(x + x2) * 0.5, (y1 + y) * 0.5, radius * 0.5}
				chreq <- Position{(x + x2) * 0.5, (y + y2) * 0.5, radius * 0.5}
			}()
		}
	}
}

func Run(tokenFile string, outfile *os.File) {
	rand.Seed(time.Now().UnixNano())
	tokens = LoadToken(tokenFile)
	fout = outfile

	const R float64 = 1000.0
	const K float64 = R * 0.001

	pos := []Position{}
	for x := 120.85000000; x < 122.2; x += 0.008984 * K {
		for y := 30.6666666; y < 31.88333333; y += 0.008984 * K {
			pos = append(pos, Position{x, y, 0})
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
			chreq <- Position{p.lng, p.lat, R}
		}
		time.Sleep(time.Second * 10)
		log.Println("Round finished")
	}
}
