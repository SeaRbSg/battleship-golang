package common

type Turn struct {
	GameID   int      `json:"game_id"`
	Response Response `json:"response"`
	Guess    Guess    `json:"guess"`
}

type Response struct {
	Hit  bool   `json:"hit"`
	Sunk string `json:"sunk,omitempty"`
	Lost bool   `json:"lost"`
}

type Guess struct {
	Guess string `json:"guess"`
}

type NewGameResponse struct {
	ID int `json:"game_id"`
}
