package podio

import (
	"flag"
	"testing"
	"time"
)

var (
	appId        int
	appSecret    string
	clientId     string
	clientSecret string
)

func init() {
	flag.IntVar(&appId, "app.id", 0, "The app id to use for tests")
	flag.StringVar(&appSecret, "app.secret", "", "The app secret to use for tests")
	flag.StringVar(&clientId, "client.id", "", "The client secret to use for tests")
	flag.StringVar(&clientSecret, "client.secret", "", "The client secret to use for tests")
}

// To use these test:
//  - Create a Podio api key. This will be used as clientId and clientSecret.
//  - Create a test app in a sandbox workspace. Get the app Id and app Secret and use them here.
//  - The test app must contain:
//    - A text field as the first field (title field)
//    - A single-choice, non-required category field (category labels are not important)
// - Run with `go test -v -app.id $APP_ID -app.secret $APP_SECREET -client.id $CLIENT_ID -client.secret $CLIENT_SECRET`

func getPodioAuth(t *testing.T) *AuthToken {
	if appId == 0 || appSecret == "" || clientId == "" || clientSecret == "" {
		t.Skip("skipping tests since no client id or secret has been set")
	}

	a, err := AuthWithAppCredentials(clientId, clientSecret, int64(appId), appSecret)
	if err != nil {
		t.Fatalf("could not generate auth token for tests: %s", err)
	}

	return a
}

func TestItemCreateAndDelete(t *testing.T) {
	// This is a quick test that we can create and delete items.

	c := NewClient(getPodioAuth(t))

	i1, err := c.CreateItem(appId, "test123", map[string]interface{}{
		"title":    "podio-go test case " + time.Now().String(),
		"category": 1,
	})

	if err != nil {
		t.Fatal("could not create item", err)
	}

	i2, err := c.CreateItem(appId, "test123", map[string]interface{}{
		"title":    "podio-go test case " + time.Now().String(),
		"category": 2,
	})
	if err != nil {
		t.Error("could not create item", err)
	}

	t.Log("created items", i1, i2)

	err = c.UpdateItem(i1, map[string]interface{}{
		"relationship": i2,
	})

	if err != nil {
		t.Error("could not update item", err)
	}

	err = c.DeleteItem(i1)
	if err != nil {
		t.Error("could not delete item", i1)
	}

	err = c.DeleteItem(i2)
	if err != nil {
		t.Error("could not delete item", i2)
	}

}

func TestItemFieldValues(t *testing.T) {
	// Tests that Field values can be unmarshalled and put into their proper types.

	c := NewClient(getPodioAuth(t))

	i, err := c.CreateItem(appId, "ab123", map[string]interface{}{
		"category": 1,
	})
	if err != nil {
		t.Fatal("could not create item", err)
	}
	defer func() {
		err = c.DeleteItem(i)
		if err != nil {
			t.Fatal("could not delete item again", i)
		}
	}()

	item, err := c.GetItem(i)
	if err != nil {
		t.Fatal("err getting item from podio", i, err)
	}

	for _, field := range item.Fields {
		switch field.Type {
		case "category":
			val, ok := field.Values.([]CategoryValue)
			if !ok {
				t.Fatalf("Expected category field value to have type %T, but found %T (value %v)",
					[]CategoryValue{},
					field.Values,
					field.Values)
			}

			if len(val) != 1 {
				t.Errorf("only expected 1 entry in the category values slice, found %d (contents: %#v)", len(val), val)
			}

			if val[0].Value.Id != 1 {
				t.Errorf("Expected category to be set to option id 1 - found, %v (full value %#v)", val[0].Value.Id, val[0].Value)
			}
		}
	}

}
