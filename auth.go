package podio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthToken struct {
	AccessToken   string                 `json:"access_token"`
	TokenType     string                 `json:"token_type"`
	ExpiresIn     int                    `json:"expires_in"`
	RefreshToken  string                 `json:"refresh_token"`
	Ref           map[string]interface{} `json:"ref"`
	TransferToken string                 `json:"transfer_token"`
}

func AuthWithUserCredentials(clientId string, clientSecret string, username string, password string) (*AuthToken, error) {
	data := url.Values{
		"grant_type":    {"password"},
		"username":      {username},
		"password":      {password},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	return authRequest(data)
}

func AuthWithAppCredentials(clientId, clientSecret string, appId uint, appToken string) (*AuthToken, error) {
	data := url.Values{
		"grant_type":    {"app"},
		"app_id":        {fmt.Sprintf("%d", appId)},
		"app_token":     {appToken},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	return authRequest(data)
}

func AuthWithAuthCode(clientId, clientSecret, authCode, redirectUri string) (*AuthToken, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectUri},
		"code":          {authCode},
	}

	return authRequest(data)
}

func authRequest(data url.Values) (*AuthToken, error) {
	var authToken AuthToken

	resp, err := http.PostForm("https://api.podio.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		podioErr := &Error{}
		err := json.Unmarshal(respBody, podioErr)
		if err != nil {
			return nil, errors.New(string(respBody))
		}
		return nil, podioErr
	}

	err = json.Unmarshal(respBody, &authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}
