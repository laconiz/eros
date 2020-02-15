package main

import (
	"encoding/json"
	"log"
)

type Struct struct {
	B struct {
		C int
	} `json:"-"`
	D json.RawMessage
}

func main() {

	s := &Struct{
		B: struct {
			C int
		}{C: 10},
		D: nil,
	}

	s.D, _ = json.Marshal(s.B)
	raw, _ := json.Marshal(s)
	log.Println(string(raw))
}
