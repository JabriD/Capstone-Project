package main

type StandingsResponse struct {
	League struct {
		Standard struct {
			SeasonYear    int `json:"seasonYear"`
			SeasonStageID int `json:"seasonStageId"`
			Conference    struct {
				East []struct {
					TeamID                 string `json:"teamId"`
					TeamSitesOnly          struct {
						TeamNickname       string `json:"teamNickname"`
					} `json:"teamSitesOnly"`
				} `json:"east"`
				West []struct {
					TeamID                 string `json:"teamId"`
					TeamSitesOnly          struct {
						TeamNickname       string `json:"teamNickname"`
					} `json:"teamSitesOnly"`
				} `json:"west"`
			} `json:"conference"`
		} `json:"standard"`
	} `json:"league"`
}

type Player struct {
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
}
type AllPlayers struct {
	League struct {
		Standard []Player `json:"standard"`
	} `json:"league"`
}