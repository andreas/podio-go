package podio

import (
	"encoding/json"
	"fmt"
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
	ExternalId         string   `json:"external_id"`
}

// partialField is used for JSON unmarshalling
type partialField struct {
	Id         int64  `json:"field_id"`
	ExternalId string `json:"external_id"`
	Type       string `json:"type"`
	Label      string `json:"label"`
	Config     struct {
		Description   string          `json:"description"`
		Required      bool            `json:"required"`
		HiddenIfEmpty bool            `json:"hidden"`
		AlwaysHidden  bool            `json:"hidden_create_view_edit"`
		ConfigVersion int             `json:"delta"`
		SettingsJSON  json.RawMessage `json:"settings"`
		Settings      interface{}     `json:"-"`
	} `json:"config"`
	ValuesJSON json.RawMessage `json:"values"`
}

// Field describes a Podio field object
type Field struct {
	partialField
	Values interface{}
}

func (f *Field) unmarshalInto(val, settings interface{}) error {
	if err := json.Unmarshal(f.ValuesJSON, val); err != nil {
		return fmt.Errorf("cannot unmarshal %q into %T: %v", f.ValuesJSON, val, err)
	}

	// allow for tests that does not set field configs.
	if settings != nil && len(f.Config.SettingsJSON) > 0 {
		if err := json.Unmarshal(f.Config.SettingsJSON, settings); err != nil {
			return fmt.Errorf("cannot unmarshal %q into %T: %v", f.Config.SettingsJSON, settings, err)
		}
	}

	return nil
}

func (f *Field) UnmarshalJSON(data []byte) error {
	f.partialField = partialField{}
	if err := json.Unmarshal(data, &f.partialField); err != nil {
		return err
	}
	var err error

	switch f.Type {
	case "app":
		values, cfg := []AppValue{}, AppFieldSettings{}
		err = f.unmarshalInto(&values, &cfg)
		f.Values, f.Config.Settings = values, cfg
	case "date":
		values := []DateValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "text":
		values, cfg := []TextValue{}, TextFieldSettings{}
		err = f.unmarshalInto(&values, &cfg)
		f.Values, f.Config.Settings = values, cfg
	case "number":
		values := []NumberValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "image":
		values := []ImageValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "member":
		values := []MemberValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "contact":
		values := []ContactValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "money":
		values := []MoneyValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "progress":
		values := []ProgressValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "location":
		values := []LocationValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "video":
		values := []VideoValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "duration":
		values := []DurationValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "embed":
		values := []EmbedValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "question":
		values := []QuestionValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "category":
		values, cfg := []CategoryValue{}, CategoryFieldSettings{}
		err = f.unmarshalInto(&values, &cfg)
		f.Values, f.Config.Settings = values, cfg
	case "tel":
		values := []TelValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	case "calculation":
		values := []CalculationValue{}
		err = f.unmarshalInto(&values, nil)
		f.Values = values
	default:
		// Unknown field type
		values, cfg := []interface{}{}, map[string]interface{}{}
		err = f.unmarshalInto(&values, &cfg)
		f.Values = values
	}
	if err != nil {
		return err
	}

	f.ValuesJSON = nil
	return nil
}

// TextValue is the value for fields of type `text`
type TextValue struct {
	Value string `json:"value"`
}

// TextFieldSettings is the configuration of a text field
type TextFieldSettings struct {
	Format string `json:"format"`
	Size   string `json:"size"`
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
	Start *Time `json:"start_utc"`
	End   *Time `json:"end_utc"`
}

// AppValue is the value for fields of type `app`
type AppValue struct {
	Value Item `json:"value"`
}

// AppFieldSettings configures an app field
type AppFieldSettings struct {
	Mulitple       bool `json:"multiple"`
	ReferencedApps []struct {
		ViewId int `json:"view_id"`
		AppId  int `json:"app_id"`
		App    App `json:"app"`
	} `json:"referenced_apps"`
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

// CategoryOption is a possible value for a category field
type CategoryOption struct {
	Status string `json:"status"`
	Text   string `json:"text"`
	Id     int    `json:"id"`
	Color  string `json:"color"`
}

// CategoryValue is the value for fields of type `category`
type CategoryValue struct {
	Value CategoryOption `json:"value"`
}

// CategoryFieldSettings holds the configuration of category fields, along with
// the possible values for the category field.
type CategoryFieldSettings struct {
	Multiple bool   `json:"multiple"`
	Display  string `json:"display"`
	Options  []CategoryOption
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
