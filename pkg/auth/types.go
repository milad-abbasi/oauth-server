package auth

type RegisterInfo struct {
	Name     string
	Email    string
	Password string
}

type LoginInfo struct {
	Email    string
	Password string
}

type Tokens struct {
	TokenType    string  `json:"token_type"`
	ExpiresIn    float64 `json:"expires_in"`
	Scope        string  `json:"scope"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
}
