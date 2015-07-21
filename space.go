package podio

import "fmt"

type Space struct {
	Id       int64  `json:"space_id"`
	Slug     string `json:"url_label"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	URLLabel string `json:"url_label"`
	OrgId    int64  `json:"org_id"`
	Push     Push   `json:"push"`
}

func (client *Client) GetSpaces(orgId int64) (spaces []Space, err error) {
	path := fmt.Sprintf("/org/%d/space", orgId)
	err = client.Request("GET", path, nil, nil, &spaces)
	return
}

func (client *Client) GetSpace(id int64) (space *Space, err error) {
	path := fmt.Sprintf("/space/%d", id)
	err = client.Request("GET", path, nil, nil, &space)
	return
}

func (client *Client) GetSpaceByOrgIdAndSlug(orgId int64, slug string) (space *Space, err error) {
	path := fmt.Sprintf("/space/org/%d/%s", orgId, slug)
	err = client.Request("GET", path, nil, nil, &space)
	return
}
