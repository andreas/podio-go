package podio

import "fmt"

// https://developers.podio.com/doc/notifications/mark-notifications-as-viewed-by-ref-553653
func (client *Client) NotificationMarkAsViewedForRef(refType string, refId int64) (statusCode int, err error) {
	path := fmt.Sprintf("/notification/%s/%d/viewed", refType, refId)
	err = client.Request("POST", path, nil, nil, &statusCode)
	return
}
