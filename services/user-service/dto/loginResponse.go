package dto

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
}
