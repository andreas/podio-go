package podio

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type File struct {
	Id   uint   `json:"file_id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Size int    `json:"size"`
}

func (client *Client) GetFiles() (files []File, err error) {
	err = client.request("GET", "/file", nil, nil, &files)
	return
}

func (client *Client) GetFile(file_id uint) (file *File, err error) {
	err = client.request("GET", fmt.Sprintf("/file/%d", file_id), nil, nil, &file)
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

func (client *Client) ReplaceFile(oldFileId, newFileId uint) error {
	path := fmt.Sprintf("/file/%d/replace", newFileId)
	params := map[string]interface{}{
		"old_file_id": oldFileId,
	}

	return client.requestWithParams("POST", path, nil, params, nil)
}

func (client *Client) AttachFile(fileId uint, refType string, refId uint) error {
	path := fmt.Sprintf("/file/%d/attach", fileId)
	params := map[string]interface{}{
		"ref_type": refType,
		"ref_id":   refId,
	}

	return client.requestWithParams("POST", path, nil, params, nil)
}

func (client *Client) DeleteFile(fileId uint) error {
	path := fmt.Sprintf("/file/%d", fileId)
	return client.request("DELETE", path, nil, nil, nil)
}
