Client for [Podio](https://podio.com) written in Go.

## Example

```go
func main() {
  authToken, err := podio.AuthWithUserCredentials("my-client-id", "my-client-secret", "my-username", "my-password")

  if err != nil {
    fmt.Println("Auth failed: ", err)
    return
  }

  client := podio.NewClient(authToken)
  orgs, err := client.GetOrganizations()

  if err != nil {
    fmt.Println("Failed to get orgs: ", err)
    return
  }

  for _, org := range orgs {
    fmt.Println("Org: ", org.Name)
  }
}
```

See [example/main.go](example/main.go).

## Item Field Values

The values of an item field depend on the type of the field. As such, the type of `Field.Values` is `interface{}` and must be coerced to access. The mapping from field types to values are as follows:

- `app`: `[]AppValue`
- `date`: `[]DateValue`
- `text`: `[]TextValue`
- `number`: `[]NumberValue`
- `image`: `[]ImageValue`
- `member`: `[]MemberValue`
- `contact`: `[]ContactValue`
- `money`: `[]MoneyValue`
- `progress`: `[]ProgressValue`
- `location`: `[]LocationValue`
- `video`: `[]VideoValue`
- `duration`: `[]DurationValue`
- `embed`: `[]EmbedValue`
- `question`: `[]QuestionValue`
- `category`: `[]CategoryValue`
- `tel`: `[]TelValue`
- `calculation`: `[]CalculationValue`

Coercing `Field.Values` safely can be done with a `switch` on `Field.Type` using the above mapping, or a type switch on `Field.Values`:

```go
// Example 1
switch field.Type {
case "app":
  for _, appVal := range field.Values.([]AppValue) {
    fmt.Println(appVal.Value.AppItemId)
  }
case "date":
  // etc
}

// Example 2
switch values := field.Values.(type) {
case []AppValue:
	for _, appVal := range values {
		fmt.Println(values[0].Value.AppItemId)
	}
case []DateValue:
  // etc
}
```

## Status

- The client supports authentication with username and password (see [Username and Password flow](https://developers.podio.com/authentication/username_password)), app authentication (see [App authentication flow](https://developers.podio.com/authentication/app_auth)) and server-side flow (see [Server-side flow](https://developers.podio.com/authentication/server_side)).
- Only supports a fraction of the API methods available, specifically around organizations, spaces, apps, items and files.
- Only a few number of fields have been defined per type.

## Contributors

The following people have contributed to podio-go:

- andreas
- stengaard
