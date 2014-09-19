package podio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	authToken  *AuthToken
}

// A podio ID is ID's for resources issued by Podio.
type ID int64

type Organization struct {
	Id   ID     `json:"org_id"`
	Slug string `json:"url_label"`
	Name string `json:"name"`
}

type Space struct {
	Id   ID     `json:"space_id"`
	Slug string `json:"url_label"`
	Name string `json:"name"`
}

type App struct {
	Id   ID     `json:"app_id"`
	Name string `json:"name"`
}

type Item struct {
	Id                 ID      `json:"item_id"`
	AppItemId          uint    `json:"app_item_id"`
	FormattedAppItemId string  `json:"app_item_id_formatted"`
	Title              string  `json:"title"`
	Files              []File  `json:"files"`
	Fields             []Field `json:"fields"`
	client             *Client
}

type Field struct {
	ExternalID string  `json:"external_id"`
	Values     []Value `json:"values"`

	Type    string `json:"type"`
	Label   string `json:"label"`
	FieldID ID     `json:"field_id"`
}

type Value struct {
	Value interface{} `json:"value"`
}

type ItemList struct {
	Filtered uint   `json:"filtered"`
	Total    uint   `json:"total"`
	Items    []Item `json:"items"`
}

type File struct {
	Id   ID     `json:"file_id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Size int    `json:"size"`
}

type Comment struct {
	ID         ID                     `json:"comment_id"`
	Value      string                 `json:"value"`
	Ref        map[string]interface{} `json:"ref"`
	Files      []File                 `json:"files"`
	CreatedBy  interface{}            `json:"created_by"`
	CreatedVia interface{}            `json:"created_via"`
	CreatedOn  interface{}            `json:"created_on"`
	IsLiked    bool                   `json:"is_liked"`
	LikeCount  int                    `json:"like_count"`
}

type AuthToken struct {
	AccessToken   string                 `json:"access_token"`
	TokenType     string                 `json:"token_type"`
	ExpiresIn     int                    `json:"expires_in"`
	RefreshToken  string                 `json:"refresh_token"`
	Ref           map[string]interface{} `json:"ref"`
	TransferToken string                 `json:"transfer_token"`
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

func AuthWithAppCredentials(client_id, client_secret string, app_id ID, app_token string) (*AuthToken, error) {
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

func NewClient(authToken *AuthToken) *Client {
	return &Client{
		httpClient: &http.Client{},
		authToken:  authToken,
	}
}

func (client *Client) request(method string, path string, headers map[string]string, body io.Reader, out interface{}) error {
	if client == nil {
		return errors.New("client not ready")
	}

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
		p := &Error{}
		err := json.Unmarshal(respBody, p)
		if err != nil {
			return errors.New(string(respBody))
		}
		return p
	}

	if out != nil {
		err = json.Unmarshal(respBody, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) GetOrganizations() (orgs []Organization, err error) {
	err = client.request("GET", "/org", nil, nil, &orgs)
	return
}

func (client *Client) GetOrganization(id ID) (org *Organization, err error) {
	path := fmt.Sprintf("/org/%d", id)
	err = client.request("GET", path, nil, nil, &org)
	return
}

func (client *Client) GetOrganizationBySlug(slug string) (org *Organization, err error) {
	path := fmt.Sprintf("/org/url?org_slug=%s", slug)
	err = client.request("GET", path, nil, nil, &org)
	return
}

func (client *Client) GetSpaces(org_id ID) (spaces []Space, err error) {
	path := fmt.Sprintf("/org/%d/space", org_id)
	err = client.request("GET", path, nil, nil, &spaces)
	return
}

func (client *Client) GetSpace(id ID) (space *Space, err error) {
	path := fmt.Sprintf("/space/%d", id)
	err = client.request("GET", path, nil, nil, &space)
	return
}

func (client *Client) GetSpaceByOrgIdAndSlug(org_id ID, slug string) (space *Space, err error) {
	path := fmt.Sprintf("/space/org/%d/%s", org_id, slug)
	err = client.request("GET", path, nil, nil, &space)
	return
}

func (client *Client) GetApps(space_id ID) (apps []App, err error) {
	path := fmt.Sprintf("/app/space/%d?view=micro", space_id)
	err = client.request("GET", path, nil, nil, &apps)
	return
}

func (client *Client) GetApp(id ID) (app *App, err error) {
	path := fmt.Sprintf("/app/%d?view=micro", id)
	err = client.request("GET", path, nil, nil, &app)
	return
}

func (client *Client) GetAppBySpaceIdAndSlug(space_id ID, slug string) (app *App, err error) {
	path := fmt.Sprintf("/app/space/%d/%s", space_id, slug)
	err = client.request("GET", path, nil, nil, &app)
	return
}

func (client *Client) GetItems(app_id ID) (items *ItemList, err error) {
	path := fmt.Sprintf("/item/app/%d/filter?fields=items.fields(files)", app_id)
	err = client.request("POST", path, nil, nil, &items)
	return
}

func (client *Client) GetItemByAppItemId(app_id ID, formatted_app_item_id string) (item *Item, err error) {
	path := fmt.Sprintf("/app/%d/item/%s", app_id, formatted_app_item_id)
	err = client.request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItemByExternalID(app_id ID, external_id string) (item *Item, err error) {
	path := fmt.Sprintf("/item/app/%d/external_id/%s", app_id, external_id)
	err = client.request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItem(item_id ID) (item *Item, err error) {
	path := fmt.Sprintf("/item/%d?fields=files", item_id)
	err = client.request("GET", path, nil, nil, &item)
	return
}

func (client *Client) CreateItem(app_id ID, external_id string, fieldValues map[string]interface{}) (item_id ID, err error) {
	path := fmt.Sprintf("/item/app/%d", app_id)
	val := map[string]interface{}{
		"fields": fieldValues,
	}

	if external_id != "" {
		val["external_id"] = external_id
	}

	buf, err := json.Marshal(val)
	if err != nil {
		return
	}

	s := &struct {
		ItemId ID `json:"item_id"`
	}{}
	err = client.request("POST", path, nil, bytes.NewReader(buf), s)
	item_id = s.ItemId

	return
}

func (client *Client) UpdateItem(itemId ID, fieldValues map[string]interface{}) error {
	path := fmt.Sprintf("/item/%d", itemId)
	buf, err := json.Marshal(map[string]interface{}{"fields": fieldValues})
	if err != nil {
		return err
	}

	return client.request("PUT", path, nil, bytes.NewBuffer(buf), nil)

}

func (client *Client) CommentOnItem(item_id ID, comment string) (err error) {
	return client.Comment("item", item_id, comment)
}

func (client *Client) Comment(typ string, objId ID, comment string) (err error) {
	path := fmt.Sprintf("/comment/%s/%d/", typ, objId)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{comment})
	if err != nil {
		return err
	}

	f := map[string]interface{}{}

	err = client.request("POST", path, nil, bytes.NewReader(buf), &f)
	return

}

func (client *Client) GetComments(typ string, objId string) (comments []*Comment, err error) {
	path := fmt.Sprintf("/comment/%s/%s/", typ, objId)
	err = client.request("GET", path, nil, nil, &comments)
	return
}

func (client *Client) GetFiles() (files []File, err error) {
	err = client.request("GET", "/file", nil, nil, &files)
	return
}

func (client *Client) GetFile(file_id string) (file *File, err error) {
	err = client.request("GET", fmt.Sprintf("/file/%s", file_id), nil, nil, &file)
	return
}

func (client *Client) GetFileContents(url string) ([]byte, error) {
	link := fmt.Sprintf("%s?oauth_token=%s", url, client.authToken.AccessToken)
	resp, err := http.Get(link)

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (client *Client) CreateFile(name string, contents []byte) (file *File, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("source", name)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(contents)
	if err != nil {
		return nil, err
	}

	err = writer.WriteField("filename", name)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	err = client.request("POST", "/file", headers, body, &file)
	return
}

func (client *Client) ReplaceFile(oldFileId, newFileId ID) error {
	path := fmt.Sprintf("/file/%d/replace", newFileId)
	body := strings.NewReader(fmt.Sprintf("{\"old_file_id\":%d}", oldFileId))
	return client.request("POST", path, nil, body, nil)
}

func (client *Client) AttachFile(fileId ID, refType, refId string) error {
	path := fmt.Sprintf("/file/%d/attach", fileId)
	body := strings.NewReader(fmt.Sprintf("{\"ref_type\":\"%s\",\"ref_id\":%s}", refType, refId))
	return client.request("POST", path, nil, body, nil)
}

func (client *Client) DeleteFile(fileId ID) error {
	path := fmt.Sprintf("/file/%d", fileId)
	return client.request("DELETE", path, nil, nil, nil)
}
