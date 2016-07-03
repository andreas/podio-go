package podio

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Item describes a Podio item object
type Item struct {
	Id                 int64    `json:"item_id"`
	AppItemId          int      `json:"app_item_id"`
	FormattedAppItemId string   `json:"app_item_id_formatted"`
	Title              string   `json:"title"`
	Files              []*File  `json:"files"`
	Fields             []*Field `json:"fields"`
	Space              Space    `json:"space"`
	App                App      `json:"app"`
	CreatedVia         Via      `json:"created_via"`
	CreatedBy          ByLine   `json:"by_line"`
	CreatedOn          Time     `json:"created_on"`
	Link               string   `json:"link"`
	Revision           int      `json:"revision"`
	Push               Push     `json:"push"`
}

// partialField is used for JSON unmarshalling
type partialField struct {
	Id         int64           `json:"field_id"`
	ExternalId string          `json:"external_id"`
	Type       string          `json:"type"`
	Label      string          `json:"label"`
	ValuesJSON json.RawMessage `json:"values"`
}

// Field describes a Podio field object
type Field struct {
	partialField
	Values interface{}
}

func (f *Field) UnmarshalJSON(data []byte) error {
	pField := partialField{}
	if err := json.Unmarshal(data, &pField); err != nil {
		return err
	}

	switch pField.Type {
	case "app":
		f.Values = []AppValue{}
	case "date":
		f.Values = []DateValue{}
	case "text":
		f.Values = []TextValue{}
	case "number":
		f.Values = []NumberValue{}
	case "image":
		f.Values = []ImageValue{}
	case "member":
		f.Values = []MemberValue{}
	case "contact":
		f.Values = []ContactValue{}
	case "money":
		f.Values = []MoneyValue{}
	case "progress":
		f.Values = []ProgressValue{}
	case "location":
		f.Values = []LocationValue{}
	case "video":
		f.Values = []VideoValue{}
	case "duration":
		f.Values = []DurationValue{}
	case "embed":
		f.Values = []EmbedValue{}
	case "question":
		f.Values = []QuestionValue{}
	case "category":
		f.Values = []CategoryValue{}
	case "tel":
		f.Values = []TelValue{}
	case "calculation":
		f.Values = []CalculationValue{}
	default:
		// Unknown field type
		f.Values = []interface{}{}
	}

	if err := json.Unmarshal(pField.ValuesJSON, &f.Values); err != nil {
		return fmt.Errorf("[ERR] Cannot unmarshal %s into %s: %v\n", string(pField.ValuesJSON), reflect.TypeOf(f.Values), err)
	}

	pField.ValuesJSON = nil
	f.partialField = pField
	return nil
}

// TextValue is the value for fields of type `text`
type TextValue struct {
	Value string `json:"value"`
}

// NumberValue is the value for fields of type `number`
type NumberValue struct {
	Value float64 `json:"value,string"`
}

// Image is the value for fields of type `image`
type ImageValue struct {
	Value File `json:"value"`
}

// DateValue is the value for fields of type `date`
type DateValue struct {
	Start *Time `json:"start"`
	End   *Time `json:"end"`
}

// AppValue is the value for fields of type `app`
type AppValue struct {
	Value Item `json:"value"`
}

// MemberValue is the value for fields of type `member`
type MemberValue struct {
	Value int `json:"value"`
}

// ContactValue is the value for fields of type `contact`
type ContactValue struct {
	Value Contact `json:"value"`
}

// MoneyValue is the value for fields of type `money`
type MoneyValue struct {
	Value    float64 `json:"value,string"`
	Currency string  `json:"currency"`
}

// ProgressValue is the value for fields of type `progress`
type ProgressValue struct {
	Value int `json:"value"`
}

// LocationValue is the value for fields of type `location`
type LocationValue struct {
	Value        string `json:"value"`
	Formatted    string `json:"formatted"`
	StreetNumber string `json:"street_number"`
	StreetName   string `json:"street_name"`
	PostalCode   string `json:"postal_code"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
}

// VideoValue is the value for fields of type `video`
type VideoValue struct {
	Value int `json:"value"`
}

// DurationValue is the value for fields of type `duration`
type DurationValue struct {
	Value int `json:"value"`
}

// EmbedValue is the value for fields of type `embed`
type EmbedValue struct {
	Embed Embed `json:"embed"`
	File  File  `json:"file"`
}

// CategoryValue is the value for fields of type `category`
type CategoryValue struct {
	Value struct {
		Status string `json:"status"`
		Text   string `json:"text"`
		Id     int    `json:"id"`
		Color  string `json:"color"`
	} `json:"value"`
}

// QuestionValue is the value for fields of type `question`
type QuestionValue struct {
	Value int `json:"value"`
}

// TelValue is the value for fields of type `tel`
type TelValue struct {
	Value int    `json:"value"`
	URI   string `json:"uri"`
}

// CalcationValue is the value for fields of type `calculation` (currently untyped)
type CalculationValue map[string]interface{}

type ItemList struct {
	Filtered int     `json:"filtered"`
	Total    int     `json:"total"`
	Items    []*Item `json:"items"`
}

func (client *Client) GetItems(appId int64) (items *ItemList, err error) {
	path := fmt.Sprintf("/item/app/%d/filter?fields=items.fields(files)", appId)
	err = client.Request("POST", path, nil, nil, &items)
	return
}

// https://developers.podio.com/doc/items/export-items-4235696
func (client *Client) ExportItems(appId int, exportFormat string, params map[string]interface{}) (int64, error) {
	path := fmt.Sprintf("/item/app/%d/export/%s", appId, exportFormat)
	rsp := &struct {
		BatchId int64 `json:"batch_id"`
	}{}

	err := client.RequestWithParams("POST", path, nil, params, rsp)

	return rsp.BatchId, err
}

func (client *Client) GetItemByAppItemId(appId int64, formattedAppItemId string) (item *Item, err error) {
	path := fmt.Sprintf("/app/%d/item/%s", appId, formattedAppItemId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItemByExternalID(appId int64, externalId string) (item *Item, err error) {
	path := fmt.Sprintf("/item/app/%d/external_id/%s", appId, externalId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) GetItem(itemId int64) (item *Item, err error) {
	path := fmt.Sprintf("/item/%d?fields=files", itemId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

func (client *Client) CreateItem(appId int, externalId string, fieldValues map[string]interface{}) (int64, error) {
	path := fmt.Sprintf("/item/app/%d", appId)
	params := map[string]interface{}{
		"fields": fieldValues,
	}

	if externalId != "" {
		params["external_id"] = externalId
	}

	rsp := &struct {
		ItemId int64 `json:"item_id"`
	}{}
	err := client.RequestWithParams("POST", path, nil, params, rsp)

	return rsp.ItemId, err
}

func (client *Client) UpdateItem(itemId int, fieldValues map[string]interface{}) error {
	path := fmt.Sprintf("/item/%d", itemId)
	params := map[string]interface{}{
		"fields": fieldValues,
	}

	return client.RequestWithParams("PUT", path, nil, params, nil)
}
