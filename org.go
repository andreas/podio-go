package podio

import "fmt"

type Organization struct {
	Id   uint   `json:"org_id"`
	Slug string `json:"url_label"`
	Name string `json:"name"`
}

func (client *Client) GetOrganizations() (orgs []Organization, err error) {
	err = client.request("GET", "/org", nil, nil, &orgs)
	return
}

func (client *Client) GetOrganization(id uint) (org *Organization, err error) {
	path := fmt.Sprintf("/org/%d", id)
	err = client.request("GET", path, nil, nil, &org)
	return
}

func (client *Client) GetOrganizationBySlug(slug string) (org *Organization, err error) {
	path := fmt.Sprintf("/org/url?org_slug=%s", slug)
	err = client.request("GET", path, nil, nil, &org)
	return
}
