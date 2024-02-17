package entity

type OAuthManager interface {
	GetOauthAccessToken() string
	GetUserInfoByOauthToken(accessToken string) *map[string]string
	CreateLinkForOAuthToken() string
}
