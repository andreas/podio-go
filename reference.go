package podio

// Reference is a reference to from one object to another Podio object
type Reference struct {
	Id       int                    `json:"id"`
	Type     string                 `json:"type"`
	TypeName string                 `json:"type_name"`
	Title    string                 `json:"title"`
	Link     string                 `json:"link"`
	Data     map[string]interface{} `json:"data"`

	CreatedOn  Time   `json:"created_on"`
	CreatedBy  ByLine `json:"created_by"`
	CreatedVia Via    `json:"created_via"`
}
