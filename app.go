package podio

import "fmt"

type App struct {
	Id              int64  `json:"app_id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	DefaultViewId   int    `json:"default_view_id"`
	URLAdd          string `json:"url_add"`
	IconId          int    `json:"icon_id"`
	LinkAdd         string `json:"link_add"`
	CurrentRevision int    `json:"current_revision"`
	ItemName        string `json:"item_name"`
	Link            string `json:"link"`
	URL             string `json:"url"`
	URLLabel        string `json:"url_label"`
	SpaceId         int    `json:"space_id"`
	Icon            string `json:"icon"`
}

// https://developers.podio.com/doc/applications/get-apps-by-space-22478
func (client *Client) GetApps(spaceId int64) (apps []App, err error) {
	path := fmt.Sprintf("/app/space/%d?view=micro", spaceId)
	err = client.Request("GET", path, nil, nil, &apps)
	return
}

// https://developers.podio.com/doc/applications/get-app-22349
func (client *Client) GetApp(id int64) (app *App, err error) {
	path := fmt.Sprintf("/app/%d?view=micro", id)
	err = client.Request("GET", path, nil, nil, &app)
	return
}

// https://developers.podio.com/doc/applications/get-app-on-space-by-url-label-477105
func (client *Client) GetAppBySpaceIdAndSlug(spaceId int64, slug string) (app *App, err error) {
	path := fmt.Sprintf("/app/space/%d/%s", spaceId, slug)
	err = client.Request("GET", path, nil, nil, &app)
	return
}
