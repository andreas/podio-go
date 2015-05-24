package podio

import "fmt"

type Space struct {
	Id       int    `json:"space_id"`
	Slug     string `json:"url_label"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	URLLabel string `json:"url_label"`
	OrgId    int    `json:"org_id"`
}

func (client *Client) GetSpaces(orgId int) (spaces []Space, err error) {
	path := fmt.Sprintf("/org/%d/space", orgId)
	err = client.Request("GET", path, nil, nil, &spaces)
	return
}

func (client *Client) GetSpace(id int) (space *Space, err error) {
	path := fmt.Sprintf("/space/%d", id)
	err = client.Request("GET", path, nil, nil, &space)
	return
}

func (client *Client) GetSpaceByOrgIdAndSlug(orgId int, slug string) (space *Space, err error) {
	path := fmt.Sprintf("/space/org/%d/%s", orgId, slug)
	err = client.Request("GET", path, nil, nil, &space)
	return
}
