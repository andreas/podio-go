package podio

import (
	"encoding/json"
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

func AuthWithUserCredentials(client_id string, client_secret string, username string, password string) (*AuthToken, error) {
	var authToken AuthToken

	data := url.Values{
		"grant_type":    {"password"},
		"username":      {username},
		"password":      {password},
		"client_id":     {client_id},
		"client_secret": {client_secret},
	}
	resp, err := http.PostForm("https://api.podio.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}

func AuthWithAppCredentials(client_id, client_secret string, app_id uint, app_token string) (*AuthToken, error) {
	var authToken AuthToken

	data := url.Values{
		"grant_type":    {"app"},
		"app_id":        {fmt.Sprintf("%d", app_id)},
		"app_token":     {app_token},
		"client_id":     {client_id},
		"client_secret": {client_secret},
	}
	resp, err := http.PostForm("https://api.podio.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}
