package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	//"os"

	_ "github.com/mattn/go-sqlite3"
)
var standResp StandingsResponse
var playerResp AllPlayers

type Seed struct {
	db *sql.DB
}

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

func createRoster() {
	//os.Remove("roster.db")

	db, err := sql.Open("sqlite3", "roster.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists teams (teamid string primary key not null, teamname text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	for _, teams := range standResp.League.Standard.Conference.East{
		stmt, _ := db.Prepare(`INSERT INTO teams(teamid, teamname) VALUES (?, ?)`)

		_, err := stmt.Exec(teams.TeamID, teams.TeamSitesOnly.TeamNickname)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seed(s Seed, createRoster string) {
	// Get the reflect value of the method
	m := reflect.ValueOf(s).MethodByName(createRoster)
	// Exit if the method doesn't exist
	if !m.IsValid() {
		log.Fatal("No method called ", createRoster)
	}
	// Execute the method
	log.Println("Seeding", createRoster, "...")
	m.Call(nil)
	log.Println("Seed", createRoster, "succedd")
}

func Execute(db *sql.DB, createRoster ...string) {
	s := Seed{db}

	seedType := reflect.TypeOf(s)

	// Execute all seeders if no method name is given
	if len(createRoster) == 0 {
		log.Println("Running all seeder...")
		// We are looping over the method on a Seed struct
		for i := 0; i < seedType.NumMethod(); i++ {
			// Get the method in the current iteration
			method := seedType.Method(i)
			// Execute seeder
			seed(s, method.Name)
		}
	}

	// Execute only the given method names
	for _, roster := range createRoster {
		seed(s, roster)
	}
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local&multiStatements=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
			// connect DB
			db, err := sql.Open("roster.db", connString)
			if err != nil {
				log.Fatalf("Error opening DB: %v", err)
			}
			Execute(db, args[1:]...)
			os.Exit(0)
		}
	}
}


func main() {
	getLeagueData()
	getPlayerData()
	createRoster()
	handleArgs()
	
	

	/*
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
	*/

}