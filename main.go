package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Response struct {
	Data []Data `json:"data"`
}

type Data struct {
	Role  string `json:"role"`
	Stats struct {
		Games int `json:"games"`
		Runes struct {
			Build []int
		} `json:"most_common_runes"`
		RuneShards struct {
			Build []int
		} `json:"rune_stat_shards"`
		Spells struct {
			Build []int `json:"build"`
		} `json:"spells"`
		Build struct {
			Build []int
		} `json:"most_common_big_item_builds"`
		StartingItems struct {
			Build []int `json:"build"`
		} `json:"most_common_starting_items"`
	} `json:"stats"`
}

type RunePayload struct {
	Name            string `json:"name"`
	PrimaryStyleId  int    `json:"primaryStyleId"`
	SelectedPerkIds []int  `json:"selectedPerkIds"`
	SubStyleId      int    `json:"subStyleId"`
}

type Q struct {
	GameData struct {
		Queue struct {
			Id int
		} `json:"queue"`
	} `json:"gameData"`
}

type Gameflow struct {
	Phase string `json:"phase"`
}

var Localhost, Auth, Champion, Queue string

func main() {
	Init()
	Update()

	loop := true
	var response Response
	champhis := "0"

	fmt.Println("Waiting for champion select.")

	for loop {
		for !champReady() {
			time.Sleep(time.Second)
		}

		Champion, _ = ApiGET("/lol-champ-select/v1/current-champion")
		Queue = getQType()

		response = getResponse()
		payload := preparePayload(response)

		json, err := json.Marshal(payload)
		Check(err)

		if champhis != Champion {
			ApiDELETE("/lol-perks/v1/pages")
			ApiPOST("/lol-perks/v1/pages", json)
			champhis = Champion
		}

		time.Sleep(time.Second)
		loop = !loopCheck()
	}
	printInfo(response)
}

func loopCheck() bool {
	body, statuscode := ApiGET("/lol-gameflow/v1/session")
	var phase Gameflow
	if statuscode != 200 {
		return false
	}
	json.Unmarshal([]byte(body), &phase)
	if phase.Phase == "InProgress" {
		return true
	} else {
		return false
	}
}

func role(data []Data) Data {
	index := 0
	max := data[0].Stats.Games

	for i, e := range data {
		if e.Stats.Games > max {
			max = e.Stats.Games
			index = i
		}
	}

	return data[index]
}

func champReady() bool {
	body, statuscode := ApiGET("/lol-champ-select/v1/current-champion")
	if statuscode != 200 || body == "0" {
		return false
	} else {
		return true
	}
}

func getResponse() Response {
	var returnresp Response
	url := "https://league-champion-aggregate.iesdev.com/api/champions/" + Champion + "?patch=" + Currentpatch + "&queue=" + Queue + "&region=world&tier=PLATINUM_PLUS"
	body := Get(url)
	err := json.Unmarshal([]byte(body), &returnresp)
	Check(err)
	return returnresp
}

func preparePayload(response Response) RunePayload {
	mostPlayed := role(response.Data).Stats
	payload := RunePayload{
		Name:            ChampList[Champion],
		PrimaryStyleId:  mostPlayed.Runes.Build[0],
		SelectedPerkIds: []int{mostPlayed.Runes.Build[1], mostPlayed.Runes.Build[2], mostPlayed.Runes.Build[3], mostPlayed.Runes.Build[4], mostPlayed.Runes.Build[6], mostPlayed.Runes.Build[7], mostPlayed.RuneShards.Build[0], mostPlayed.RuneShards.Build[1], mostPlayed.RuneShards.Build[2]},
		SubStyleId:      mostPlayed.Runes.Build[5],
	}
	return payload
}

func getQType() string {
	var q Q
	body, _ := ApiGET("/lol-gameflow/v1/session")
	json.Unmarshal([]byte(body), &q)
	if q.GameData.Queue.Id == 450 {
		return "450"
	} else {
		return "420"
	}
}

func printInfo(resp Response) {
	mostPlayed := role(resp.Data)
	spells := mostPlayed.Stats.Spells.Build
	buildsq := mostPlayed.Stats.Build.Build
	finalbuild := []string{ItemList[buildsq[0]], ItemList[buildsq[1]], ItemList[buildsq[2]], ItemList[buildsq[3]], ItemList[buildsq[4]], ItemList[buildsq[5]]}
	buildsq = mostPlayed.Stats.StartingItems.Build
	startingitems := []string{ItemList[buildsq[0]], ItemList[buildsq[1]]}

	fmt.Println("\n========================================================================")
	fmt.Println("Champion:", ChampList[Champion], "\tRole:", mostPlayed.Role)
	fmt.Println("Summoners:", SummonerSpells[spells[0]], SummonerSpells[spells[1]])
	fmt.Println("\nStarting Items:\n", strings.Join(startingitems, " | "))
	fmt.Println("Final Build:\n", strings.Join(finalbuild, " | "))
	fmt.Println("\n========================================================================")
}
