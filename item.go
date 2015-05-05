package podio

import "fmt"

type Item struct {
	Id                 uint     `json:"item_id"`
	AppItemId          uint     `json:"app_item_id"`
	FormattedAppItemId string   `json:"app_item_id_formatted"`
	Title              string   `json:"title"`
	Files              []*File  `json:"files"`
	Fields             []*Field `json:"fields"`
}

type Field struct {
	FieldId    uint     `json:"field_id"`
	ExternalId string   `json:"external_id"`
	Type       string   `json:"type"`
	Label      string   `json:"label"`
	Values     []*Value `json:"values"`
}

type Value struct {
	Value interface{} `json:"value"`
}

type ItemList struct {
	Filtered uint    `json:"filtered"`
	Total    uint    `json:"total"`
	Items    []*Item `json:"items"`
}

func (client *Client) GetItems(app_id uint) (items *ItemList, err error) {
	path := fmt.Sprintf("/item/app/%d/filter?fields=items.fields(files)", app_id)
	err = client.Request("POST", path, nil, nil, &items)
	return
}

func (client *Client) GetItemByAppItemId(app_id uint, formatted_app_item_id string) (item *Item, err error) {
	path := fmt.Sprintf("/app/%d/item/%s", app_id, formatted_app_item_id)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItemByExternalID(app_id uint, external_id string) (item *Item, err error) {
	path := fmt.Sprintf("/item/app/%d/external_id/%s", app_id, external_id)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItem(item_id uint) (item *Item, err error) {
	path := fmt.Sprintf("/item/%d?fields=files", item_id)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) CreateItem(app_id uint, external_id string, fieldValues map[string]interface{}) (uint, error) {
	path := fmt.Sprintf("/item/app/%d", app_id)
	params := map[string]interface{}{
		"fields": fieldValues,
	}

	if external_id != "" {
		params["external_id"] = external_id
	}

	rsp := &struct {
		ItemId uint `json:"item_id"`
	}{}
	err := client.RequestWithParams("POST", path, nil, params, rsp)

	return rsp.ItemId, err
}

func (client *Client) UpdateItem(itemId uint, fieldValues map[string]interface{}) error {
	path := fmt.Sprintf("/item/%d", itemId)
	params := map[string]interface{}{
		"fields": fieldValues,
	}

	return client.RequestWithParams("PUT", path, nil, params, nil)
}
