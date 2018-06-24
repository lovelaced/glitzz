package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Boards struct {
	board string
}

func getBoards(msgs chan<- string) {
	println("getting boards...")
	resp, err := http.Get("https://a.4cdn.org/boards.json")
	if err != nil {
		println("Failed to get 4chan boards.json")
	}
	body, err := ioutil.ReadAll(resp.Body)
	jsonbody := string(body)
	if err != nil {
		println("Can't parse boards.json body")
	}
	parsedBody := boards{}
	jerr := json.Unmarshal([]byte(jsonbody), &parsedBody)
	if jerr != nil {
		println("fucked up the unmarshal")
	}
	println(Boards.board)

}
