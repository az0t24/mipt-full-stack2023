package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OAuthManager struct {
	ClientID     string
	ClientSecret string
}

func NewOAuthManager(client_id, client_secret string) *OAuthManager {
	return &OAuthManager{
		ClientID:     client_id,
		ClientSecret: client_secret,
	}
}

type YandexRes struct {
	Email     string `json:"default_email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (o *OAuthManager) GetUserInfoByOauthToken(accessToken string) *map[string]string {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://login.yandex.ru/info?format=json", nil)
	req.Header.Add("Authorization", "OAuth "+accessToken)
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Printf("Response body: %v\n", string(body))

	userOauth := YandexRes{}
	json.Unmarshal(body, &userOauth)
	log.Printf("userRegOauth: %+v\n", userOauth)

	oauthMap := convertYandexResToMap(&userOauth)

	return &oauthMap
}

func (o *OAuthManager) GetOauthAccessToken() string {
	// TODO: fix this boilerplate accessToken
	accessToken := "y0_AgAEA7qkNgu_AArQYQAAAADx8Qs09Lyllt5KTviIbTssc7p0cfGl2Qw"

	return accessToken
}

func convertYandexResToMap(yandexRes *YandexRes) map[string]string {
	email := yandexRes.Email
	firstName := yandexRes.FirstName
	lastName := yandexRes.LastName
	res := map[string]string{
		"Email":     email,
		"FirstName": firstName,
		"LastName":  lastName,
	}
	return res
}

func (o *OAuthManager) CreateLinkForOAuthToken() string {
	reqURL := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=token&client_id=%s", o.ClientID)
	return reqURL
}
