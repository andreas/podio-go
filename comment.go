package podio

import "fmt"

// Comment is a comment on an object in podio.
// The object to which this comment is associated is described in this Reference.
type Comment struct {
	Id         uint       `json:"comment_id"`
	ExternalId string     `json:"external_id"`
	Value      string     `json:"value"`
	Ref        *Reference `json:"ref"`
	Files      []*File    `json:"files"`
	CreatedBy  ByLine     `json:"created_by"`
	CreatedVia Via        `json:"created_via"`
	CreatedOn  Time       `json:"created_on"`
	IsLiked    bool       `json:"is_liked"`
	LikeCount  int        `json:"like_count"`
}

// CommentAttr adds a comment to a podio object. It returns a Comment (with podio ID) or an error if one occured.
//
// refType (item, task, ...) and refId identifies the podio object to which the comment is added.
// Additional parameters can be set in the params map.
// text is the actual comment value.
func (client *Client) CommentAttr(refType string, refId uint, params map[string]interface{}, text string) (*Comment, error) {
	path := fmt.Sprintf("/comment/%s/%d/", refType, refId)
	if params == nil {
		params = map[string]interface{}{}
	}
	params["value"] = text

	comment := &Comment{}
	err := client.requestWithParams("POST", path, nil, params, comment)
	return comment, err
}

// Comment adds a comment to a podio object.
//
// This is equivalent to calling CommentAttr with a nil params map.
func (client *Client) Comment(refType string, refId uint, text string) (*Comment, error) {
	return client.CommentAttr(refType, refId, nil, text)
}

// GetComments retrieves the comments associated with a podio object.
//
// refType is the type of the podio object. For legal type values see
// refId is the podio id of the podio object.
func (client *Client) GetComments(refType string, refId uint) (comments []*Comment, err error) {
	path := fmt.Sprintf("/comment/%s/%d/", refType, refId)
	err = client.request("GET", path, nil, nil, &comments)
	return
}
