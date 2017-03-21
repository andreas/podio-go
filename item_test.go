package podio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"reflect"

	"github.com/kr/pretty"
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

// To add more tests using the fixtures simply generate the situation you are trying to
// reproduce in Podio, fetch the results from https://developers.podio.com/ and add the JSON
// data to the fixtures directory. This sidesteps the authentication dance in testing.

// TestUnmarshalFieldValuesAndSettings tests that all values in a given Podio item can be
// unmarshalled correctly. The fixture json files are simply the output from
// https://developers.podio.com/doc/items/get-item-22360 with appropriate item IDs
func TestUnmarshalFieldValuesAndSettings(t *testing.T) {

	errOnUnknownField = true
	defer func() {
		errOnUnknownField = false
	}()

	type fv struct {
		ignore   string
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
					settings: ContactFieldSettings{"space_users", []string{"user"}},
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
		{
			jsonFile: "fixtures/item_350017179.json",
			fields: []fv{
				{
					values:   []TextValue{{"incom"}},
					settings: TextFieldSettings{"plain", "small"},
				},
				{
					values:   []TextValue{{"<p>sadfas</p>"}},
					settings: TextFieldSettings{"html", "large"},
				},
				{
					values:   []PhoneValue{{"+19197936102", "work"}},
					settings: PhoneFieldSettings{"callto", []string{"mobile", "work", "home", "main", "work_fax", "private_fax", "other"}},
				},
				{
					values:   []EmailValue{{"hhah@podio.com", "work"}},
					settings: EmailFieldSettings{PossibleTypes: []string{"work", "home", "other"}},
				},
				{
					values: []CategoryValue{{Value: CategoryOption{Status: "active", Text: "C", Id: 3, Color: "DCEBD8"}}},
					settings: CategoryFieldSettings{
						Display: "inline",
						Options: []CategoryOption{
							{"active", "A", 1, "DCEBD8"},
							{"active", "B", 2, "DCEBD8"},
							{"active", "C", 3, "DCEBD8"},
						},
					},
				},
				{
					values:   []TextValue{{"1012312312313123123123123123123123123123123"}},
					settings: TextFieldSettings{"plain", "small"},
				},
				{
					values: []CalculationValue{{"111.0000"}},
					settings: CalculationFieldSettings{
						Script:     "var a = @[Number](field_80608151);\nparseInt(a, 10);",
						ReturnType: "number",
						Decimals:   2,
					},
				},
				{
					values: []DateValue{{
						Start: parseTime(t, "2016-06-24 19:00:00"),
						End:   parseTime(t, "2016-06-30 20:00:00"),
					}},
					settings: DateFieldSettings{Calendar: true, End: "enabled", Time: "enabled"},
				},
				{ignore: "same as the test above. Repeat field."},
				{
					values: []AppValue{{Value: Item{
						Id: 331300398, AppItemId: 3, Title: "sadf 41`", Revision: 14,
						Link:  "https://podio.com/podio/sandbox-4fcx2i/apps/allfields/items/3",
						Files: []*File{},
						Space: Space{
							Id:       2720177,
							Name:     "sandbox",
							URL:      "https://podio.com/podio/sandbox-4fcx2i",
							URLLabel: "sandbox-4fcx2i",
							OrgId:    736,
						},
						App: App{
							Id:              10421272,
							Name:            "AllFields",
							Status:          "active",
							Link:            "https://podio.com/podio/sandbox-4fcx2i/apps/allfields",
							URL:             "https://podio.com/podio/sandbox-4fcx2i/apps/allfields",
							URLAdd:          "https://podio.com/podio/sandbox-4fcx2i/apps/allfields/items/new",
							LinkAdd:         "https://podio.com/podio/sandbox-4fcx2i/apps/allfields/items/new",
							CurrentRevision: 25,
							ItemName:        "Merge",
							URLLabel:        "allfields",
							SpaceId:         2720177,
							IconId:          251,
							Icon:            "251.png",
						},
						CreatedVia: Via{Id: 1, Name: "Podio"},
						CreatedOn:  *parseTime(t, "2015-10-09 12:49:43"),
					}}},
					settings: AppFieldSettings{
						Mulitple: true,
						ReferencedApps: []struct {
							AppId  int `json:"app_id"`
							App    App `json:"app"`
							ViewId int `json:"view_id"`
						}{{
							AppId: 10421272,
							App: App{
								Id:              10421272,
								Name:            "AllFields",
								Status:          "active",
								Link:            "https://podio.com/podio/sandbox-4fcx2i/apps/allfields",
								URL:             "https://podio.com/podio/sandbox-4fcx2i/apps/allfields",
								URLAdd:          "https://podio.com/podio/sandbox-4fcx2i/apps/allfields/items/new",
								LinkAdd:         "https://podio.com/podio/sandbox-4fcx2i/apps/allfields/items/new",
								CurrentRevision: 25,
								ItemName:        "Merge",
								URLLabel:        "allfields",
								SpaceId:         2720177,
								IconId:          251,
								Icon:            "251.png",
							},
							ViewId: 0,
						}, {
							AppId: 10118971,
							App: App{
								Id:     10118971,
								Name:   "GPG Keys",
								Status: "active", URLAdd: "https://podio.com/podio/sandbox-4fcx2i/apps/gpg-keys/items/new",
								IconId:          14,
								LinkAdd:         "https://podio.com/podio/sandbox-4fcx2i/apps/gpg-keys/items/new",
								CurrentRevision: 6,
								ItemName:        "GPG Key",
								Link:            "https://podio.com/podio/sandbox-4fcx2i/apps/gpg-keys",
								URL:             "https://podio.com/podio/sandbox-4fcx2i/apps/gpg-keys",
								URLLabel:        "gpg-keys",
								SpaceId:         2720177,
								Icon:            "14.png",
							},
							ViewId: 0,
						}},
					},
				},
				{ignore: "Slight variaion of calculation field. Still with number output"},
				{
					values: []ContactValue{{Value: Contact{
						UserId:  2468975,
						SpaceId: 0,
						Type:    "user",
						Image: File{
							Id:   297530466,
							Link: "https://d2cmuesa4snpwn.cloudfront.net/public/297530466",
						},
						ProfileId:  140798621,
						Link:       "https://podio.com/users/2468975",
						Avatar:     297530466,
						Name:       "Brian Stengaard",
						LastSeenOn: parseTime(t, "2017-03-21 15:43:34"),
					}}},
					settings: ContactFieldSettings{"space_users", []string{"user"}},
				},
				{ignore: "we already have a test for Number field"},
				{
					values: []ImageValue{{Value: File{
						Id:   195174666,
						Name: "a1.jpeg",
						Link: "https://files.podio.com/195174666",
						Size: 45326,
					}}},
					settings: ImageFieldSettings{[]string{"image/png"}},
				},
				{ignore: "type money already covered"},
				{
					values: []EmbedValue{{
						Embed: Embed{Id: 117018345, Type: "link", Title: "Google", Description: "Search the world's information, including webpages, images, videos and more. Google has many special features to help you find exactly what you're looking for.", EmbedHTML: "", URL: "https://google.com/", OriginalURL: "https://google.com/", ResolvedURL: "https://google.com/", Hostname: "google.com", EmbedHeight: 0, EmbedWidth: 0},
						File:  File{Id: 231335541, Link: "https://files.podio.com/231335541"},
					}},
				},
			},
		},
		{
			jsonFile: "fixtures/item_582709679.json",
			fields: []fv{
				{values: []ProgressValue{{Value: 42}}},
				{
					values: []LocationValue{{
						Value:      "Champ de Mars, 5 Avenue Anatole France, 75007 Paris, France",
						Formatted:  "Champ de Mars, 5 Avenue Anatole France, 75007 Paris, France",
						PostalCode: "75007",
						City:       "Paris",
						State:      "Île-de-France",
						Country:    "France",
						Lat:        48.8583701,
						Lng:        2.2944813}},
					settings: LocationFieldSettings{false, true},
				},
				{values: []DurationValue{{93784}}, settings: DurationFieldSettings{[]string{"days", "hours", "minutes", "seconds"}}},
			},
		},
	}

	for _, c := range testCases {
		t.Run("field values test of "+c.jsonFile, func(t *testing.T) {
			t.Logf("Checking field values on %q", c.jsonFile)

			item := &Item{}
			err := json.Unmarshal(getFixtureJSON(t, c.jsonFile), item)
			if err != nil {
				t.Errorf("Could not unmarshal item: %s", err)
				return
			}

			if len(c.fields) != len(item.Fields) {
				t.Errorf("expected %d fields, but found %d", len(c.fields), len(item.Fields))
			}
			for i, f := range item.Fields {
				if i >= len(c.fields) {
					break
				}
				if c.fields[i].ignore != "" {
					continue
				}

				if !reflect.DeepEqual(f.Values, c.fields[i].values) {
					t.Errorf("mismatching values on field num %d: %s (type: %s field_id: %d)\nGot: %# v\nExp: %# v\nDiff: %v\n",
						i+1, f.Label, f.Type, f.Id,
						pretty.Formatter(f.Values),
						pretty.Formatter(c.fields[i].values),
						pretty.Diff(f.Values, c.fields[i].values))
				}

				if !reflect.DeepEqual(f.Config.Settings, c.fields[i].settings) {
					t.Errorf("mismatching settings on field num %d: %s (type: %s field_id: %d)\nGot: %#v\nExp: %#v\nDiff: %v\n",
						i+1, f.Label, f.Type, f.Id,
						pretty.Formatter(f.Config.Settings),
						pretty.Formatter(c.fields[i].settings),
						pretty.Diff(f.Config.Settings, c.fields[i].settings))
				}
			}
		})
	}
}

// parseTime is a test helper to parse the UTC timestamps found various places in Podio data
func parseTime(t *testing.T, tim string) *Time {
	out := &Time{}

	err := out.UnmarshalJSON([]byte(fmt.Sprintf("%q", tim)))
	if err != nil {
		t.Fatalf("could not parse time %q as a Podio time: %v", tim, err)
	}
	return out
}

// getFixtureJSON reads a fixture file and returns the raw buffer. It aborts the
// test framework if an error occurs while reading the file.
func getFixtureJSON(t *testing.T, name string) []byte {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatalf("Could not fetch fixture file: %v\n", err)
	}

	return buf
}
