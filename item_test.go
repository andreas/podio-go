package podio

import (
	"encoding/json"
	"testing"

	"io/ioutil"

	"reflect"

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

func TestUnmarshalItem(t *testing.T) {
	item := &Item{}
	err := json.Unmarshal(getFixtureJSON(t, "item"), item)
	if err != nil {
		t.Fatal("Could not unmarshal item")
	}

	testCases := []struct {
		values   interface{}
		settings interface{}
	}{
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
	}

	if len(testCases) != len(item.Fields) {
		t.Errorf("expected %d fields, but found %d", len(testCases), len(item.Fields))
	}

	for i, f := range item.Fields {
		if i >= len(testCases) {
			break
		}

		if !reflect.DeepEqual(f.Values, testCases[i].values) {
			t.Errorf("mismatching values on field num %d: %s (%s)\nGot: %+v\nExp: %+v\n",
				i+1, f.Label, f.Type, f.Values, testCases[i].values)
		}
		if !reflect.DeepEqual(f.Config.Settings, testCases[i].settings) {
			t.Errorf("mismatching settings on field num %d: %s (%s)\nGot: %+v\nExp: %+v\n",
				i+1, f.Label, f.Type, f.Config.Settings, testCases[i].settings)
		}

	}
}

func getFixtureJSON(t *testing.T, name string) []byte {
	fname := "fixtures/" + name + ".json"
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("Could not fetch fixture file: %v\n", err)
	}

	return buf
}
