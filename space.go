package podio

import "fmt"

type Space struct {
	Id   uint   `json:"space_id"`
	Slug string `json:"url_label"`
	Name string `json:"name"`
}

func (client *Client) GetSpaces(org_id uint) (spaces []Space, err error) {
	path := fmt.Sprintf("/org/%d/space", org_id)
	err = client.request("GET", path, nil, nil, &spaces)
	return
}

func (client *Client) GetSpace(id uint) (space *Space, err error) {
	path := fmt.Sprintf("/space/%d", id)
	err = client.request("GET", path, nil, nil, &space)
	return
}

func (client *Client) GetSpaceByOrgIdAndSlug(org_id uint, slug string) (space *Space, err error) {
	path := fmt.Sprintf("/space/org/%d/%s", org_id, slug)
	err = client.request("GET", path, nil, nil, &space)
	return
}
