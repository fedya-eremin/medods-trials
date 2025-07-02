package dto

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Lifetime     int    `json:"lifetime"`
}

type LoginResponse TokenPair
