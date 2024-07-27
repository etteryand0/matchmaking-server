package users

type User struct {
	UserId      string   `json:"user_id"`
	MMR         int      `json:"mmr"`
	Roles       []string `json:"roles"`
	WaitingTime int      `json:"waitingTime"`
}
