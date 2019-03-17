package podio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	authToken  *AuthToken
}

type Error struct {
	Parameters interface{} `json:"error_parameters"`
	Detail     interface{} `json:"error_detail"`
	Propagate  bool        `json:"error_propagate"`
	Request    struct {
		URL   string `json:"url"`
		Query string `json:"query_string"`
	} `json:"request"`
	Description string `json:"error_description"`
	Type        string `json:"error"`
}

func (p *Error) Error() string {
	return fmt.Sprintf("%s: %s", p.Type, p.Description)
}

func NewClient(authToken *AuthToken) *Client {
	return &Client{
		httpClient: &http.Client{},
		authToken:  authToken,
	}
}

func (client *Client) Request(method string, path string, headers map[string]string, body io.Reader, out interface{}) error {
	req, err := http.NewRequest(method, "https://api.podio.com"+path, body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Authorization", "OAuth2 "+client.authToken.AccessToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		podioErr := &Error{}
		err := json.Unmarshal(respBody, podioErr)
		if err != nil {
			return errors.New(string(respBody))
		}
		return podioErr
	}

	if out != nil {
		return json.Unmarshal(respBody, out)
	}

	return nil
}

func (client *Client) RequestWithParams(method string, path string, headers map[string]string, params map[string]interface{}, out interface{}) error {
    var body io.Reader

    if method == "GET" {
        pathURL, err := url.Parse(path)
        if err != nil {
            return err
        }
        query := pathURL.Query()
        for key, value := range params {
            query.Add(key, fmt.Sprint(value))
        }
        pathURL.RawQuery = query.Encode()
        path = pathURL.String()
    } else {
        buf, err := json.Marshal(params)
        if err != nil {
            return err
        }
        body = bytes.NewReader(buf)
    }

    return client.Request(method, path, headers, body, out)
}
