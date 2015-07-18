package podio

// Contact describes a Podio contact object
type Contact struct {
	UserId     int    `json:"user_id"`
	SpaceId    int    `json:"space_id"`
	Type       string `json:"type"`
	Image      File   `json:"image"`
	ProfileId  int    `json:"profile_id"`
	OrgId      int    `json:"org_id"`
	Link       string `json:"link"`
	Avatar     int    `json:"avatar"`
	LastSeenOn *Time  `json:"last_seen_on"`
	Name       string `json:"name"`
}
