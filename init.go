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

var Currentpatch string
var SummonerSpells map[int]string = map[int]string{21: "Barrier", 1: "Cleanse", 14: "Ignite", 3: "Exhaust", 4: "Flash", 6: "Ghost", 7: "Heal", 13: "Clarity", 11: "Smite", 39: "Mark", 32: "Mark", 12: "Teleport"}

func Init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	Localhost, Auth = GetVars("C:/Riot Games/League of Legends/lockfile")

	var versions []string

	url := "https://ddragon.leagueoflegends.com/api/versions.json"
	body := Get(url)
	err := json.Unmarshal([]byte(body), &versions)
	Check(err)
	Currentpatch = versions[0]
}

func GetVars(location string) (string, string) {
	file, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatal("Run League Of Legends before SimplyRunes")
	}
	array := strings.Split(string(file), ":")
	return fmt.Sprintf("https://127.0.0.1:%s", array[2]), "Basic " + toauth(array[3]) // 2 = port, 3 = pw
}

func toauth(pass string) string {
	convertthis := "riot:" + pass
	return base64.StdEncoding.EncodeToString([]byte(convertthis))
}
