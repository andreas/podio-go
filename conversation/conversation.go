package conversation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/andreas/podio-go"
)

type Client struct {
	*podio.Client
}

// Metadata holds meta-data about a group or direct chat session
type Metadata struct {
	ConversationId uint            `json:"conversation_id"`
	Reference      podio.Reference `json:"ref"`
	CreatedOn      podio.Time      `json:"created_on"`
	CreatedBy      podio.ByLine    `json:"created_by"`

	Excerpt      string         `json:"excerpt"`
	Starred      bool           `json:"starred"`
	Unread       bool           `json:"unread"`
	UnreadCount  uint           `json:"unread_count"`
	LastEvent    podio.Time     `json:"last_event_on"`
	Subject      string         `json:"subject"`
	Participants []podio.ByLine `json:"participants"`
	Type         string         `json:"type"` // direct or group
}

// ConversationEvent is a single message from a sender to a conversation
type Event struct {
	EventID uint `json:"event_id"`

	Action string `json:"action"`
	Data   struct {
		MessageID uint          `json:"message_id"`
		Files     []interface{} `json:"files"` // TODO: add structure
		Text      string        `json:"text"`
		EmbedFile interface{}   `json:"embed_file"` // TODO: add structure
		Embed     interface{}   `json:"embed"`      // TODO: add structure
		CreatedOn podio.Time
	}

	CreatedVia podio.Via    `json:"created_via"`
	CreatedBy  podio.ByLine `json:"created_by"`
	CreatedOn  podio.Time   `json:"created_on"`
}

// ConversationSelector can modify the scope of a conversations lookup request - see WithLimit and WithOffset for examples.
type Selector func(uri *url.URL)

// GetConversation returns all conversations that the client has access to (max 200). Use WithLimit and WithOffset
// to do pagination if that is what you want.
func (client *Client) GetConversations(withOpts ...Selector) ([]Metadata, error) {
	u, err := url.Parse("/conversation/")
	if err != nil { // should never happen
		return nil, err
	}
	for _, selector := range withOpts {
		selector(u)
	}

	convs := []Metadata{}
	err = client.Request("GET", u.RequestURI(), nil, nil, &convs)
	return convs, err
}

// GetConversationEvents returns all events for the conversation with id conversationId. WithLimit and WithOffset can be used to do
// pagination.
func (client *Client) GetEvents(conversationId uint, withOpts ...Selector) ([]Event, error) {
	u, err := url.Parse(fmt.Sprintf("/conversation/%d/event", conversationId))
	if err != nil { // should never happen
		return nil, err
	}
	for _, selector := range withOpts {
		selector(u)
	}

	convs := []Event{}
	err = client.Request("GET", u.RequestURI(), nil, nil, &convs)
	return convs, err
}

// Reply sends a (string) message to the conversation identified by conversationId. Only text strings are supported (that is
// no embedding for now).
func (client *Client) Reply(conversationId uint, reply string) (Event, error) {
	path := fmt.Sprintf("/conversation/%d/reply/v2", conversationId)
	out := Event{}

	buf, err := json.Marshal(map[string]string{"text": reply})
	if err != nil {
		return out, err
	}
	err = client.Request("POST", path, map[string]string{"content-type": "application/json"}, bytes.NewReader(buf), &out)
	return out, err
}

// WithLimit sets a limit on the returned list of Conversations or ConversationEvents. limit must be in the range (0-200].
func WithLimit(limit uint) Selector {
	f := func(u *url.URL) {
		q := u.Query()
		q.Add("limit", strconv.Itoa(int(limit)))
		u.RawQuery = q.Encode()
	}
	return Selector(f)
}

// WithOffset introduces an offset in the returned list of Conversations or ConversationsEvents.
func WithOffset(offset uint) Selector {
	f := func(u *url.URL) {
		q := u.Query()
		q.Add("offset", strconv.Itoa(int(offset)))
		u.RawQuery = q.Encode()
	}
	return Selector(f)
}

// Unread manipulates the conversation request to only conversations with unread messages.
func Unread() Selector {
	f := func(u *url.URL) {
		u.Path = path.Join(u.Path, "unread")
	}
	return Selector(f)
}
