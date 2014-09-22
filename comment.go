package podio

import "fmt"

type Comment struct {
	ID         uint                   `json:"comment_id"`
	Value      string                 `json:"value"`
	Ref        map[string]interface{} `json:"ref"`
	Files      []*File                `json:"files"`
	CreatedBy  interface{}            `json:"created_by"`
	CreatedVia interface{}            `json:"created_via"`
	CreatedOn  interface{}            `json:"created_on"`
	IsLiked    bool                   `json:"is_liked"`
	LikeCount  int                    `json:"like_count"`
}

func (client *Client) Comment(refType, refId, text string) (*Comment, error) {
	path := fmt.Sprintf("/comment/%s/%d/", refType, refId)
	params := map[string]interface{}{
		"value": text,
	}

	comment := &Comment{}
	err := client.requestWithParams("POST", path, nil, params, comment)
	return comment, err
}

func (client *Client) GetComments(refType string, refId string) (comments []*Comment, err error) {
	path := fmt.Sprintf("/comment/%s/%s/", refType, refId)
	err = client.request("GET", path, nil, nil, &comments)
	return
}
