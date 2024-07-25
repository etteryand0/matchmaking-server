package users

type User struct {
	UserId      string   `json:"user_id"`
	MMR         uint     `json:"mmr"`
	Roles       []string `json:"roles"`
	WaitingTime uint     `json:"waitingTime"`
}
