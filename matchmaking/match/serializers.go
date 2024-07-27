package match

type TestCollection struct {
	Last string `json:"last"`
}

type Match struct {
	MatchId string `json:"match_id"`
	Teams   [2]struct {
		Side  string `json:"side"`
		Users []struct {
			ID   string `json:"id"`
			Role string `json:"role"`
		} `json:"users"`
	} `json:"teams"`
}
