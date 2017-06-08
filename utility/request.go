package utility

import (
	"log"
	"time"
	"strings"

	"io/ioutil"
	"net/http"
	"math/rand"
)

func NowMillisec() int {
	return int(time.Now().UnixNano() / 1000000)
}

func newRequest(method, url, data string) *http.Request {
	body := strings.NewReader(data)
	for {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			log.Println("http new request error:", err)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500) + 100))
		} else {
			return req
		}
	}
}

func send(client *http.Client, url, data string, setHeader func(req *http.Request)) *http.Response {
	for {
		req := newRequest("POST", url, data)
		setHeader(req)
		res, err := client.Do(req)
		if err != nil {
			log.Println("client do error:", err)
		} else {
			if res.StatusCode != 200 {
				log.Println("response status = ", res.Status)
			} else {
				return res
			}
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500) + 1000))
	}
}

func ReadAll(client *http.Client, url, data string, setHeader func(req *http.Request)) []byte {
	for {
		res := send(client, url, data, setHeader)
		defer res.Body.Close()

		text, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("read request body error:", err)
		} else {
			return text
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500) + 1000))
	}
}
