package facebook

import (
	fb "github.com/huandu/facebook"
)

var fbClient *fb.App

func InitFbClient(appId string, appSecret string) {
	fb.Version = "v1.0"

	fbClient = fb.New(appId, appSecret)
	//TODO: check for valid fb credentials after client creation
}

func GetLongTermToken(shortTermToken string) (token string, expires int, err error) {
	//TODO: check on getting actual long-term access token
	return fbClient.ExchangeToken(shortTermToken)
}
