package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type AllPlayers struct {
	League struct {
		Standard []struct {
			FirstName            string `json:"firstName"`
			LastName             string `json:"lastName"`
			TemporaryDisplayName string `json:"temporaryDisplayName,omitempty"`
			PersonID             string `json:"personId"`
			TeamID               string `json:"teamId"`
			Jersey               string `json:"jersey"`
			IsActive             bool   `json:"isActive"`
			Pos                  string `json:"pos"`
			HeightFeet           string `json:"heightFeet"`
			HeightInches         string `json:"heightInches"`
			WeightPounds         string `json:"weightPounds"`
			DateOfBirthUTC       string `json:"dateOfBirthUTC"`
			Teams []struct {
				TeamID      string `json:"teamId"`
				SeasonStart string `json:"seasonStart"`
				SeasonEnd   string `json:"seasonEnd"`
			} `json:"teams"`
			Draft struct {
				TeamID     string `json:"teamId"`
				PickNum    string `json:"pickNum"`
				RoundNum   string `json:"roundNum"`
				SeasonYear string `json:"seasonYear"`
			} `json:"draft"`
			NbaDebutYear    string `json:"nbaDebutYear"`
			YearsPro        string `json:"yearsPro"`
			CollegeName     string `json:"collegeName"`
			LastAffiliation string `json:"lastAffiliation"`
			Country         string `json:"country"`
			IsallStar       bool   `json:"isallStar,omitempty"`
		} `json:"standard"`
	} `json:"league"`
}

func main() {
	res, err := http.Get("http://data.nba.net/10s/prod/v1/2018/players.json")
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
