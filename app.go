package podio

import "fmt"

type App struct {
	Id   uint   `json:"app_id"`
	Name string `json:"name"`
}

func (client *Client) GetApps(space_id uint) (apps []App, err error) {
	path := fmt.Sprintf("/app/space/%d?view=micro", space_id)
	err = client.request("GET", path, nil, nil, &apps)
	return
}

func (client *Client) GetApp(id uint) (app *App, err error) {
	path := fmt.Sprintf("/app/%d?view=micro", id)
	err = client.request("GET", path, nil, nil, &app)
	return
}

func (client *Client) GetAppBySpaceIdAndSlug(space_id uint, slug string) (app *App, err error) {
	path := fmt.Sprintf("/app/space/%d/%s", space_id, slug)
	err = client.request("GET", path, nil, nil, &app)
	return
}
