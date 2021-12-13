package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var ChampList = make(map[string]string)
var ItemList = make(map[int]string)

type ChampData struct {
	Map map[string]struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"data"`
}

type ItemData struct {
	Map map[string]struct {
		Name string `json:"name"`
	} `json:"data"`
}

func Update() {
	url := "https://ddragon.leagueoflegends.com/cdn/" + Currentpatch + "/data/en_US/champion.json"
	body := Get(url)

	var m ChampData
	err := json.Unmarshal([]byte(body), &m)
	Check(err)

	for _, e := range m.Map {
		ChampList[e.Key] = e.Name
	}

	url = "https://ddragon.leagueoflegends.com/cdn/" + Currentpatch + "/data/en_US/item.json"
	body = Get(url)

	var n ItemData
	err = json.Unmarshal([]byte(body), &n)
	Check(err)

	for i, e := range n.Map {
		converted, err := strconv.Atoi(i)
		Check(err)
		ItemList[converted] = e.Name
	}
}

func Check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func Get(url string) string {
	req, err := http.Get(url)
	Check(err)
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	Check(err)
	return string(body)
}

// https://league-champion-aggregate.iesdev.com/apii/champons/157?patch=11.23.1&queue=420&region=world&role=MID&tier=PLATINUM_PLUS
