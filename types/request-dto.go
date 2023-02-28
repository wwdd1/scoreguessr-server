package types

type PostAuthRequestBody struct {
	IdToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
}
