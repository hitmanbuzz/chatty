package user

type LoginPayload struct {
	LoginType string `json:"type"`
	Username  string `json:"username"`
}
