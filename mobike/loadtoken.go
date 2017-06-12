package mobike

import (
	"log"
	"bytes"

	"io/ioutil"
	"encoding/json"
)

type Token struct {
	Accesstoken string `json:"accesstoken"`
	Mobileno string `json:"mobileno"`
}

func printToken(tokens []Token) {
	for _, tok := range tokens {
		log.Printf("accesstoken: %s, mobileno: %s\n", tok.Accesstoken, tok.Mobileno)
	}
}

func LoadToken(tokenFile string) []Token {
	dat, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		panic(err)
	}
	tokens := []Token{}
	err = json.Unmarshal(dat, &tokens)
	if err == nil {
		printToken(tokens)
	} else {
		for _, tok := range bytes.Split(dat, []byte("\n")) {
			if len(tok) == 0 {
				break
			}
			tmp := bytes.Split(tok, []byte(" "))
			accesstoken, mobileno := string(tmp[0]), string(tmp[1])
			tokens = append(tokens, Token{accesstoken, mobileno})
		}
		printToken(tokens)
	}

	if len(tokens) == 0 {
		panic("no tokens loaded")
	}
	return tokens
}
