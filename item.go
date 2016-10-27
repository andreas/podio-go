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

func (f *Field) unmarshalValuesInto(out interface{}) error {
	if err := json.Unmarshal(f.ValuesJSON, &out); err != nil {
		return fmt.Errorf("[ERR] Cannot unmarshal %s into %s: %v\n", f.ValuesJSON, reflect.TypeOf(out), err)
	}
	return nil
}

func (f *Field) UnmarshalJSON(data []byte) error {
	f.partialField = partialField{}
	if err := json.Unmarshal(data, &f.partialField); err != nil {
		return err
	}

	switch f.Type {
	case "app":
		values := []AppValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "date":
		values := []DateValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "text":
		values := []TextValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "number":
		values := []NumberValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "image":
		values := []ImageValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "member":
		values := []MemberValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "contact":
		values := []ContactValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "money":
		values := []MoneyValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "progress":
		values := []ProgressValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "location":
		values := []LocationValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "video":
		values := []VideoValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "duration":
		values := []DurationValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "embed":
		values := []EmbedValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "question":
		values := []QuestionValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "category":
		values := []CategoryValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "tel":
		values := []TelValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "calculation":
		values := []CalculationValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "email":
		values := []EmailValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	case "phone":
		values := []PhoneValue{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	default:
		// Unknown field type
		values := []interface{}{}
		f.unmarshalValuesInto(&values)
		f.Values = values
	}

	f.ValuesJSON = nil
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

// EmailValue is values for email fields.
type EmailValue struct {
	// Type is the use of the email. Possible values are other, home or work
	Type string `json:"type"`
	// The actual email: jill@example.com
	Value string `json:"value"`
}

// PhoneValue contains the value for phone fields.
type PhoneValue struct {
	// Type is the use of the phone field. Possible values are
	// mobile, work, home, main, work_fax, private, fax or other
	Type string `json:"type"`
	// Value contains the phone number itself
	Value string `json:"value"`
}

// CalcationValue is the value for fields of type `calculation` (currently untyped)
type CalculationValue map[string]interface{}

type ItemList struct {
	Filtered int     `json:"filtered"`
	Total    int     `json:"total"`
	Items    []*Item `json:"items"`
}

// https://developers.podio.com/doc/items/filter-items-4496747
func (client *Client) GetItems(appId int64) (items *ItemList, err error) {
	path := fmt.Sprintf("/item/app/%d/filter?fields=items.fields(files)", appId)
	err = client.Request("POST", path, nil, nil, &items)
	return
}

// https://developers.podio.com/doc/items/filter-items-4496747
func (client *Client) FilterItems(appId int64, params map[string]interface{}) (items *ItemList, err error) {
	path := fmt.Sprintf("/item/app/%d/filter?fields=items.fields(files)", appId)
	err = client.RequestWithParams("POST", path, nil, params, &items)
	return
}

// https://developers.podio.com/doc/items/get-item-by-app-item-id-66506688
func (client *Client) GetItemByAppItemId(appId int64, formattedAppItemId string) (item *Item, err error) {
	path := fmt.Sprintf("/app/%d/item/%s", appId, formattedAppItemId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

// https://developers.podio.com/doc/items/get-item-by-external-id-19556702
func (client *Client) GetItemByExternalID(appId int64, externalId string) (item *Item, err error) {
	path := fmt.Sprintf("/item/app/%d/external_id/%s", appId, externalId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

// https://developers.podio.com/doc/items/get-item-22360
func (client *Client) GetItem(itemId int64) (item *Item, err error) {
	path := fmt.Sprintf("/item/%d?fields=files", itemId)
	err = client.Request("GET", path, nil, nil, &item)
	return
}

// https://developers.podio.com/doc/items/add-new-item-22362
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

// https://developers.podio.com/doc/items/update-item-22363
func (client *Client) UpdateItem(itemId int, fieldValues map[string]interface{}) error {
	path := fmt.Sprintf("/item/%d", itemId)
	params := map[string]interface{}{
		"fields": fieldValues,
	}

	return client.RequestWithParams("PUT", path, nil, params, nil)
}
