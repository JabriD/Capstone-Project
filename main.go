package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func getPlayers() {
	res, err := http.Get("http://data.nba.net/10s/prod/v1/2022/players.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var players AllPlayers
	err = json.Unmarshal(bs, &players)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(players)
}

func getTeams () {
	res, err := http.Get("http://data.nba.net/prod/v1/current/standings_conference.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var allTeams Standings
	err = json.Unmarshal(bs, &allTeams)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(allTeams)
}

func main() {
	//getPlayers()
	getTeams()
}