package podio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalTextValue(t *testing.T) {
	r := require.New(t)

	fieldJson := []byte(`{
		"type": "text",
		"values": [{
			"value": "a"
		}]
	}`)

	field := &Field{}
	err := json.Unmarshal(fieldJson, field)
	r.NoError(err)

	values, ok := field.Values.([]TextValue)
	r.True(ok, "Expected values to be []podio.TextValue, is %#v", field.Values)
	r.Len(values, 1)
	r.Equal(values[0].Value, "a")
}

func TestUnmarshalAppValue(t *testing.T) {
	r := require.New(t)

	fieldJson := []byte(`{
		"type": "app",
		"values": [{
			"value": {
				"item_id": 1
			}
		}]
	}`)

	field := &Field{}
	err := json.Unmarshal(fieldJson, field)
	r.NoError(err)

	values, ok := field.Values.([]AppValue)
	r.True(ok, "Expected values to be []podio.AppValue, is %#v", field.Values)
	r.Len(values, 1)
	r.Equal(values[0].Value.Id, int64(1))
}

func TestUnmarshalItemValues(t *testing.T) {

	type fv struct {
		values   interface{}
		settings interface{}
	}

	testCases := []struct {
		jsonFile string
		fields   []fv
	}{
		{
			jsonFile: "fixtures/item_225607452.json",
			fields: []fv{
				{
					values:   []TextValue{{"Title"}},
					settings: TextFieldSettings{"plain", "small"},
				},
				{
					values: []CategoryValue{
						{CategoryOption{Status: "active", Text: "B", Id: 2, Color: "DCEBD8"}},
					},
					settings: CategoryFieldSettings{Display: "inline", Options: []CategoryOption{
						{Status: "active", Text: "A", Id: 1, Color: "DCEBD8"},
						{Status: "active", Text: "B", Id: 2, Color: "DCEBD8"},
						{Status: "active", Text: "C", Id: 3, Color: "DCEBD8"},
					}},
				},
				{
					values: []DateValue{{
						Start: parseTime(t, "2014-12-11 22:00:00"),
					}},
					settings: DateFieldSettings{
						Calendar: true, End: "enabled", Time: "enabled",
					},
				},
				{
					values: []ContactValue{{Value: Contact{
						UserId: 2468975,
						Type:   "user",
						Image: File{
							Id:   125807791,
							Link: "https://d2cmuesa4snpwn.cloudfront.net/public/125807791",
						},
						ProfileId:  140798621,
						Link:       "https://podio.com/users/2468975",
						Avatar:     125807791,
						Name:       "Brian Stengaard",
						LastSeenOn: parseTime(t, "2014-12-10 15:26:35"),
					}}},
				},
				{
					values:   []NumberValue{{6513.5100}},
					settings: NumberFieldSettings{0},
				},
				{
					values:   []MoneyValue{{541.987, "EUR"}},
					settings: MoneyFieldSettings{[]string{"USD", "EUR"}},
				},
				{
					values: []EmbedValue{{Embed: Embed{
						Id:          55017316,
						Type:        "link",
						Title:       "Google",
						Description: "Search the world's information, including webpages, images, videos and more. Google has many special features to help you find exactly what you're looking for.",
						URL:         "http://google.com/",
						OriginalURL: "http://google.com/",
						ResolvedURL: "http://google.com/",
						Hostname:    "google.com",
					}}},
				},
			},
		},
	}

	for _, c := range testCases {
		item := &Item{}
		err := json.Unmarshal(getFixtureJSON(t, c.jsonFile), item)
		if err != nil {
			t.Errorf("Could not unmarshal item: %s", err)
			continue
		}

		if len(c.fields) != len(item.Fields) {
			t.Errorf("expected %d fields, but found %d", len(c.fields), len(item.Fields))
		}

		for i, f := range item.Fields {
			if i >= len(c.fields) {
				break
			}

			if !reflect.DeepEqual(f.Values, c.fields[i].values) {
				t.Errorf("mismatching values on field num %d: %s (%s)\nGot: %+v\nExp: %+v\n",
					i+1, f.Label, f.Type, f.Values, c.fields[i].values)
			}
			if !reflect.DeepEqual(f.Config.Settings, c.fields[i].settings) {
				t.Errorf("mismatching settings on field num %d: %s (%s)\nGot: %+v\nExp: %+v\n",
					i+1, f.Label, f.Type, f.Config.Settings, c.fields[i].settings)
			}

		}
	}
}

func parseTime(t *testing.T, tim string) *Time {
	out := &Time{}

	err := out.UnmarshalJSON([]byte(fmt.Sprintf("%q", tim)))
	if err != nil {
		t.Fatalf("could not parse time %q as a Podio time: %v", tim, err)
	}
	return out
}

func getFixtureJSON(t *testing.T, name string) []byte {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatalf("Could not fetch fixture file: %v\n", err)
	}

	return buf
}
