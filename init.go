package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const Version string = "3.0"

var CurrentPatch string
var LastPatch string
var APILink string
var Auth string

var Champs map[string]string = make(map[string]string)
var Summs map[string]string = make(map[string]string)
var Items map[string]string = make(map[string]string)

type StaticData struct {
	Data map[string]struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"data"`
}

func init() {
	fmt.Println("Loading...")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if checkForUpdates() != Version {
		fmt.Println("=============================================")
		fmt.Println("New version is available, update here:")
		fmt.Println("https://github.com/sevnnn/SimplyRunes")
		fmt.Println("=============================================")
		fmt.Println("")
	}

	getLockfile()

	getPatch()
	getChamps()
	getSumms()
	getItems()

}

func checkForUpdates() string {
	data, _ := Get("https://raw.githubusercontent.com/sevnnn/SimplyRunes/main/version")
	return string(data)
}

func getConfig() string {
	f, err := ioutil.ReadFile("./settings.json")
	Check(err)

	var settings struct {
		Path string `json:"LeaguePath"`
	}

	json.Unmarshal(f, &settings)

	path := settings.Path

	if string(path[len(path)-1]) == "/" {
		path = path[:len(path)-1]
	}

	return path
}

func getLockfile() {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/lockfile", getConfig()))
	if err != nil {
		log.Fatalln("Run League Of Legends before running SimplyRunes")
	}

	array := strings.Split(string(file), ":")
	APILink = fmt.Sprintf("%s://127.0.0.1:%s", array[4], array[2])
	Auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("riot:"+array[3])))
}

func getPatch() {
	var patchlist []string
	data, _ := Get("https://ddragon.leagueoflegends.com/api/versions.json")
	json.Unmarshal(data, &patchlist)

	CurrentPatch = patchlist[0]
	LastPatch = patchlist[1]
}

func getChamps() {
	var sd StaticData
	data, _ := Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", CurrentPatch))
	json.Unmarshal(data, &sd)

	for _, e := range sd.Data {
		Champs[e.Key] = e.Name
	}
}

func getSumms() {
	var sd StaticData
	data, _ := Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/summoner.json", CurrentPatch))
	json.Unmarshal(data, &sd)

	for _, e := range sd.Data {
		Summs[e.Key] = e.Name
	}
}

func getItems() {
	var sd StaticData
	data, _ := Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/item.json", CurrentPatch))
	json.Unmarshal(data, &sd)

	for i, e := range sd.Data {
		Items[i] = e.Name
	}
}
