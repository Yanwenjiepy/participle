package internal

import (
	"encoding/json"
	"fmt"
	"github.com/jszwec/csvutil"
	"io/ioutil"
	"net/http"
	"participle/logger"
)

type Account struct {
	AppID     string `csv:"AppID"`
	ApiKey    string `csv:"APIKey"`
	SecretKey string `csv:"SecretKey"`
}

type accessKeyRes struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}

func LoadCSV(accountFilePath string, accounts *[]Account) error {

	accountContentByte, err := ioutil.ReadFile(accountFilePath)
	if err != nil {
		return err
	}

	err = csvutil.Unmarshal(accountContentByte, accounts)
	if err != nil {
		return err
	}

	return nil
}

func GetAccessToken(accounts *[]Account) []string {

	accessTokenList := make([]string, 0, len(*accounts))
	accessTokenRes := accessKeyRes{}

	for _, account := range *accounts {
		appID, apiKey, secretKey := account.AppID, account.ApiKey, account.SecretKey
		url := fmt.Sprintf(
			"https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
			apiKey,
			secretKey)

		response, err := http.Get(url)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to get access key of appID: %s", appID))
			continue
		}

		replyContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			_ = response.Body.Close()
			logger.Log.Error("Failed to get response body of access token request!")
			continue
		}

		_ = response.Body.Close()

		err = json.Unmarshal(replyContent, &accessTokenRes)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to unmarshal access token response and err: %s", err.Error()))
			continue
		}

		accessToken := accessTokenRes.AccessToken

		accessTokenList = append(accessTokenList, accessToken)
	}

	return accessTokenList
}
