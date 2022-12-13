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

	_ "github.com/mattn/go-sqlite3"
)

type Seed struct {
	db *sql.DB
}

func getLeagueData() (StandingsResponse, error) {
	res, err := http.Get("http://data.nba.net/prod/v1/current/standings_conference.json")
	if err != nil {
		return StandingsResponse{}, fmt.Errorf("failed to get league data: %w", err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return StandingsResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var standResp StandingsResponse
	err = json.Unmarshal(bs, &standResp)
	if err != nil {
		return StandingsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return standResp, nil
}

func getPlayerData() (AllPlayers, error) {
	res, err := http.Get("https://data.nba.net/10s/prod/v1/2022/players.json")
	if err != nil {
		return AllPlayers{}, fmt.Errorf("failed to get player data: %w", err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return AllPlayers{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var playerResp AllPlayers
	err = json.Unmarshal(bs, &playerResp)
	if err != nil {
		return AllPlayers{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return playerResp, nil
}

func createRoster(standResp StandingsResponse, playerResp AllPlayers) {
	os.Remove("roster.db")

	db, err := sql.Open("sqlite3", "roster.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists teams (teamid string primary key not null, teamname text not null);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	playerStmt := `
	create table if not exists players (teamid string not null, playername string not null, playerid string not null, foreign key(teamid) references teams(teamid));
	`
	_, err = db.Exec(playerStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, playerStmt)
		return
	}

	for _, teams := range standResp.League.Standard.Conference.East {
		stmt, _ := db.Prepare(`INSERT INTO teams(teamid, teamname) VALUES (?, ?)`)

		_, err := stmt.Exec(teams.TeamID, teams.TeamSitesOnly.TeamNickname)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, teams := range standResp.League.Standard.Conference.West {
		stmt, _ := db.Prepare(`INSERT INTO teams(teamid, teamname) VALUES (?, ?)`)

		_, err := stmt.Exec(teams.TeamID, teams.TeamSitesOnly.TeamNickname)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, players := range playerResp.League.Standard {
		stmt, _ := db.Prepare(`INSERT INTO players(teamid, playername, playerid) VALUES (?, ?, ?)`)

		_, err := stmt.Exec(players.TeamID, players.TemporaryDisplayName, players.PersonID)
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

// Example commands:
//
//	go run . -refreshdb true
//	go run . -c "SELECT * FROM players;"
func main() {
	var (
		refreshdb bool
		sqlCmd    string
	)

	// Use flag package to capture the CLI values and put them in seed and sqlCmd.
	flag.BoolVar(&refreshdb, "refreshdb", false, "Refresh Database")
	flag.StringVar(&sqlCmd, "c", "", "Database Query - Team")
	flag.Parse()

	if refreshdb {
		// debugging
		//log.Fatal("refreshing db")

		standResp, err := getLeagueData()
		if err != nil {
			log.Fatal(err)
		}

		playerResp, err := getPlayerData()
		if err != nil {
			log.Fatal(err)
		}

		createRoster(standResp, playerResp)
		log.Println("Roster Refreshed")
	}

	if sqlCmd != "" {
		playerResp, err := getPlayerData()
		if err != nil {
			log.Fatal(err)
		}

		standResp, err := getLeagueData()
		if err != nil {
			log.Fatal(err)
		}

		createRoster(standResp, playerResp)
		// debugging
		//log.Fatalf("sql CMD: %s", sqlCmd)

		// (The below can probably be improved. Errors can have more detail added, and I'm not sure if you want to just
		// print the row columns like I'm doing here. You also have to consider that users could pass malicious commands
		// to do things like delete all your data or insert/update garbage data.)

		teamQuery := "SELECT * FROM players;"

		db, err := sql.Open("sqlite3", "roster.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query(teamQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			for _, players := range playerResp.League.Standard {
				err = rows.Scan(&players.TeamID, &players.TemporaryDisplayName, &players.PersonID)
				if err != nil {
					log.Fatal(err)
				}
				if players.TeamID == sqlCmd {
					fmt.Println(players.TeamID, players.TemporaryDisplayName, players.PersonID)
				}
			}
		}
	}
	standResp, _ := getLeagueData()
	playerResp, _ := getPlayerData()
	createRoster(standResp, playerResp)
	handleArgs()
}