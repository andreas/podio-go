package podio

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type File struct {
	Id   int    `json:"file_id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Size int    `json:"size"`
}

func (client *Client) GetFiles() (files []File, err error) {
	err = client.Request("GET", "/file", nil, nil, &files)
	return
}

func (client *Client) GetFile(file_id int) (file *File, err error) {
	err = client.Request("GET", fmt.Sprintf("/file/%d", file_id), nil, nil, &file)
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

	err = client.Request("POST", "/file", headers, body, &file)
	return
}

func (client *Client) ReplaceFile(oldFileId, newFileId int) error {
	path := fmt.Sprintf("/file/%d/replace", newFileId)
	params := map[string]interface{}{
		"old_file_id": oldFileId,
	}

	return client.RequestWithParams("POST", path, nil, params, nil)
}

func (client *Client) AttachFile(fileId int, refType string, refId int) error {
	path := fmt.Sprintf("/file/%d/attach", fileId)
	params := map[string]interface{}{
		"ref_type": refType,
		"ref_id":   refId,
	}

	return client.RequestWithParams("POST", path, nil, params, nil)
}

func (client *Client) DeleteFile(fileId int) error {
	path := fmt.Sprintf("/file/%d", fileId)
	return client.Request("DELETE", path, nil, nil, nil)
}
