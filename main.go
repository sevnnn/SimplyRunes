package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type S struct {
	Path     string `json:"LeaguePath"`
	Infotype string `json:"infotype"`
}

var Settings S
var Localhost string
var Auth string

func main() {
	// get latest data
	patch, champ_id_map, spells_id_map, items_id_map := getData()

	// im bad @ programming
	showBuild := false
	var build, starting_items []int

	// load settings
	txtsettings := loadSettings()
	json.Unmarshal(txtsettings, &Settings)

	// preaparing champion history
	championhis := "0"

	// cheking if league is open
	fmt.Println("Waiting for League to run")
	for !leagueOpened() {
		time.Sleep(time.Second * 5) // cpu usage capper
	}

	// init lcu
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	Localhost, Auth = initLCU()

	// perma loop
	fmt.Println("Waiting for champ select")

	for leagueOpened() {

		// wait for champion select and checking if i should show build
		for !inChampSelect() && !showBuild {
			time.Sleep(time.Second) // cpu usage capper
		}

		// get champion
		champion, _ := ApiGET("/lol-champ-select/v1/current-champion")
		if string(champion) != championhis {
			// setting new championhis
			championhis = string(champion)

			// getting runes and build
			runes_primary, runes_secondary, runes, role, spells, bs, sis := getInfo(string(champion), patch, getQType())
			build = bs
			starting_items = sis

			// setting runes
			payload, err := json.Marshal(
				RunePayload{
					Primary:   runes_primary,
					Secondary: runes_secondary,
					Name:      champ_id_map[string(champion)],
					Runes:     runes,
				})
			Check(err)

			ApiDELETE("/lol-perks/v1/pages")
			ApiPOST("/lol-perks/v1/pages", payload)

			// printing status
			fmt.Println("Runes set for", champ_id_map[string(champion)], role, "\tSummoner spells:", spells_id_map[spells[0]], spells_id_map[spells[1]])

			// setting value
			showBuild = true
		}

		// showing build
		if inGame() && showBuild {
			// show this only one time
			showBuild = false

			// convert []int to []string that can be joined
			var realstarting_items []string
			for _, e := range starting_items {
				realstarting_items = append(realstarting_items, items_id_map[strconv.Itoa(e)])
			}
			var realbuild []string
			for _, e := range build {
				realbuild = append(realbuild, items_id_map[strconv.Itoa(e)])
			}

			// print out build and starting items
			fmt.Println("Starting items:")
			fmt.Println(strings.Join(realstarting_items, " | "))
			fmt.Println("Build:")
			fmt.Println(strings.Join(realbuild, " | "))
		}

		time.Sleep(time.Second) // cpu usage capper
	}
}

func Check(e error) {
	// error handler
	if e != nil {
		log.Fatalln(e)
	}
}

func leagueOpened() bool {
	// check if league is opened by cheking lockfile
	_, err := ioutil.ReadFile(Settings.Path + "/lockfile")
	if err != nil {
		return false
	} else {
		return true
	}
}

func loadSettings() []byte {
	// returns []byte from settings.json
	file, err := ioutil.ReadFile("./settings.json")
	Check(err)

	return file
}

func initLCU() (string, string) {
	// returning localhost and auth
	file, err := ioutil.ReadFile(Settings.Path + "/lockfile")
	Check(err)

	array := strings.Split(string(file), ":")

	return array[4] + "://127.0.0.1:" + array[2], "Basic " + base64.StdEncoding.EncodeToString([]byte("riot:"+array[3]))
}

func inChampSelect() bool {
	// checking if player is in champion select
	data, _ := ApiGET("/lol-gameflow/v1/gameflow-phase")
	if string(data) != "\"ChampSelect\"" {
		return false
	} else {
		return true
	}
}

func inGame() bool {
	data, _ := ApiGET("/lol-gameflow/v1/gameflow-phase")
	if string(data) != "\"InProgress\"" {
		return false
	} else {
		return true
	}
}

func getData() (string, map[string]string, map[string]string, map[string]string) {
	var patchlist []string
	data, _ := Get("https://ddragon.leagueoflegends.com/api/versions.json")
	json.Unmarshal(data, &patchlist)

	var champ_id_name StaticData
	var champ_id_map = make(map[string]string)
	data, _ = Get("https://ddragon.leagueoflegends.com/cdn/" + patchlist[0] + "/data/en_US/champion.json")
	json.Unmarshal(data, &champ_id_name)
	for _, e := range champ_id_name.Data {
		champ_id_map[e.Key] = e.Name
		//champ_id_map[e.Name] = e.Name
	}

	var spell_id_name StaticData
	var spell_id_map = make(map[string]string)
	data, _ = Get("https://ddragon.leagueoflegends.com/cdn/" + patchlist[0] + "/data/en_US/summoner.json")
	json.Unmarshal(data, &spell_id_name)
	for _, e := range spell_id_name.Data {
		spell_id_map[e.Key] = e.Name
		//champ_id_map[e.Name] = e.Name
	}

	var item_id_name StaticData2
	var item_id_map = make(map[string]string)
	data, _ = Get("https://ddragon.leagueoflegends.com/cdn/" + patchlist[0] + "/data/en_US/item.json")
	json.Unmarshal(data, &item_id_name)
	//fmt.Println(item_id_map)
	//fmt.Println(item_id_name)
	for i, e := range item_id_name.Data {
		item_id_map[i] = e.Name
	}

	return patchlist[0], champ_id_map, spell_id_map, item_id_map
}

func getQType() string {
	var gameflow Gameflow
	data, _ := ApiGET("/lol-gameflow/v1/session")
	json.Unmarshal(data, &gameflow)
	if gameflow.Gamedata.Queue.GameMode == "ARAM" {
		return "450"
	} else {
		return "420"
	}
}

func getInfo(champ string, patch string, qtype string) (int, int, []int, string, []string, []int, []int) {
	var apidata APIData
	data, _ := Get("https://league-champion-aggregate.iesdev.com/api/champions/" + champ + "?patch=" + patch + "&queue=" + qtype + "&region=world&tier=PLATINUM_PLUS")
	json.Unmarshal(data, &apidata)

	r := 0
	role := "ARAM"
	if qtype != "450" {
		maxgames := 0
		// determining most popular role
		for i, e := range apidata.Data {
			if e.Stats.Games > maxgames {
				maxgames = e.Stats.Games
				r = i
			}
		}
		role = apidata.Data[r].Role
	}

	if Settings.Infotype == "winrate" {
		p := apidata.Data[r].Stats.WR_Runes.Build[0]
		s := apidata.Data[r].Stats.WR_Runes.Build[5]
		rs := []int{apidata.Data[r].Stats.WR_Runes.Build[1], apidata.Data[r].Stats.WR_Runes.Build[2], apidata.Data[r].Stats.WR_Runes.Build[3], apidata.Data[r].Stats.WR_Runes.Build[4], apidata.Data[r].Stats.WR_Runes.Build[6], apidata.Data[r].Stats.WR_Runes.Build[7], apidata.Data[r].Stats.WR_Shards.Build[0], apidata.Data[r].Stats.WR_Shards.Build[1], apidata.Data[r].Stats.WR_Shards.Build[2]}
		spells := apidata.Data[r].Stats.Spells.Build
		b := apidata.Data[r].Stats.WR_Build.Build
		si := apidata.Data[r].Stats.WR_Starting_items.Build
		return p, s, rs, role, []string{strconv.Itoa(spells[0]), strconv.Itoa(spells[1])}, b, si
	} else {
		p := apidata.Data[r].Stats.Runes.Build[0]
		s := apidata.Data[r].Stats.Runes.Build[5]
		rs := []int{apidata.Data[r].Stats.Runes.Build[1], apidata.Data[r].Stats.Runes.Build[2], apidata.Data[r].Stats.Runes.Build[3], apidata.Data[r].Stats.Runes.Build[4], apidata.Data[r].Stats.Runes.Build[6], apidata.Data[r].Stats.Runes.Build[7], apidata.Data[r].Stats.Shards.Build[0], apidata.Data[r].Stats.Shards.Build[1], apidata.Data[r].Stats.Shards.Build[2]}
		b := apidata.Data[r].Stats.Build.Build
		spells := apidata.Data[r].Stats.Spells.Build
		si := apidata.Data[r].Stats.Starting_items.Build
		return p, s, rs, role, []string{strconv.Itoa(spells[0]), strconv.Itoa(spells[1])}, b, si
	}
}
