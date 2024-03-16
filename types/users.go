package types

type UserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Age    uint8  `json:"age"`
	Gender string `json:"gender"`
}

type MatchedUser struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Age    uint8  `json:"age"`
	Gender string `json:"gender"`
}
