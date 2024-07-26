package match

type MatchRespone struct {
	Epoch    string `json:"epoch"`
	TestName string `json:"test_name"`
}

type TestCollection struct {
	Last string `json:"last"`
}
