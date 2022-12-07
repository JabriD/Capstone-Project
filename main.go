package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)
var standResp StandingsResponse
var playerResp AllPlayers

func getLeagueData() (StandingsResponse, error) {
	res, err := http.Get("http://data.nba.net/prod/v1/current/standings_conference.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	err = json.Unmarshal(bs, &standResp)
	if err != nil {
		log.Fatal(err)
	}
	return standResp, nil
}

func getPlayerData() (AllPlayers, error) {
	res, err := http.Get("https://data.nba.net/10s/prod/v1/2022/players.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	err = json.Unmarshal(bs, &playerResp)
	if err != nil {
		log.Fatal(err)
	}
	return playerResp, nil
}


func main() {
	getLeagueData()
	getPlayerData()
	for _, teams := range standResp.League.Standard.Conference.East{
		fmt.Println(teams.TeamID, teams.TeamSitesOnly.TeamNickname)
	}
	for _, teams := range standResp.League.Standard.Conference.West{
		fmt.Println(teams.TeamID, teams.TeamSitesOnly.TeamNickname)
	}
	for _, players := range playerResp.League.Standard{
		if players.IsActive{
		fmt.Println(players.TeamID,players.TemporaryDisplayName, players.PersonID)
		}
	}
}