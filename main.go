package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	var champion_his string
	var show_build bool
	var build Build

	fmt.Println("Waiting for champion select")
	for {
		champ, champ_sc := ApiGet("/lol-champ-select/v1/current-champion")
		if string(champ) != "0" && champ_sc == 200 && string(champ) != champion_his && inChampSelect() {
			champion_his = string(champ)

			build = setRunesGetBuild(string(champ))
			show_build = true
		}

		if inGame() && show_build {
			show_build = false

			fmt.Printf("Starting items:\n%s\nFull Build:\n%s\n", strings.Join(build.Starting_items, " | "), strings.Join(build.Items, " | "))
		}

		time.Sleep(time.Second)
	}
}

func Check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func inChampSelect() bool {
	data, _ := ApiGet("/lol-gameflow/v1/gameflow-phase")
	if string(data) == "\"ChampSelect\"" {
		return true
	} else {
		return false
	}
}

func inGame() bool {
	data, _ := ApiGet("/lol-gameflow/v1/gameflow-phase")
	if string(data) != "\"InProgress\"" {
		return false
	} else {
		return true
	}
}

func getQType() string {
	var gameflow Gameflow
	data, _ := ApiGet("/lol-gameflow/v1/session")
	json.Unmarshal(data, &gameflow)
	if gameflow.Gamedata.Queue.GameMode == "ARAM" {
		return "450"
	} else {
		return "420"
	}
}

func setRunesGetBuild(champ string) Build {
	data, sc := Get(fmt.Sprintf("https://league-champion-aggregate.iesdev.com/api/champions/%s?patch=%s&queue=%s&region=world&tier=PLATINUM_PLUS", champ, CurrentPatch, getQType()))
	if sc != 200 {
		data, _ = Get(fmt.Sprintf("https://league-champion-aggregate.iesdev.com/api/champions/%s?patch=%s&queue=%s&region=world&tier=PLATINUM_PLUS", champ, LastPatch, getQType()))
	}

	var apidata APIData
	json.Unmarshal(data, &apidata)

	pos := 0
	role := "ARAM"
	if getQType() != "450" {
		maxgames := 0

		for i, e := range apidata.Data {
			if e.Stats.Games > maxgames {
				maxgames = e.Stats.Games
				pos = i
			}
		}
		role = apidata.Data[pos].Role
	}

	// shortcut
	short := apidata.Data[pos].Stats

	payload, err := json.Marshal(RunePayload{
		Name:      fmt.Sprintf("%s - %s", Champs[champ], role),
		Primary:   short.Runes.Build[0],
		Secondary: short.Runes.Build[5],
		Runes:     []int{short.Runes.Build[1], short.Runes.Build[2], short.Runes.Build[3], short.Runes.Build[4], short.Runes.Build[6], short.Runes.Build[7], short.Shards.Build[0], short.Shards.Build[1], short.Shards.Build[2]},
	})
	Check(err)

	ApiDELETE("/lol-perks/v1/pages")
	ApiPOST("/lol-perks/v1/pages", payload)

	fmt.Printf("\nRunes set for %s (%s)\nSummoner spells: %s %s\n\n", Champs[champ], role, Summs[strconv.Itoa(short.Spells.Build[0])], Summs[strconv.Itoa(short.Spells.Build[1])])

	var starting_items_string []string
	for _, e := range apidata.Data[pos].Stats.Starting_items.Build {
		starting_items_string = append(starting_items_string, Items[strconv.Itoa(e)])
	}

	var items_string []string
	for _, e := range apidata.Data[pos].Stats.Build.Build {
		items_string = append(items_string, Items[strconv.Itoa(e)])
	}

	return Build{
		Starting_items: starting_items_string,
		Items:          items_string,
	}
}
