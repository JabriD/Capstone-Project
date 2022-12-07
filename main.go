package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getTeams () (StandingsResponse, error) {
	res, err := http.Get("http://data.nba.net/prod/v1/current/standings_conference.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var standResp StandingsResponse
	err = json.Unmarshal(bs, &standResp)
	if err != nil {
		log.Fatal(err)
	}
	return standResp, nil
}

func main() {
	getTeams()
}