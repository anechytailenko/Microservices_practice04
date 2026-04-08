package create_meetup

type Command struct {
	Title       string `json:"title"`
	Capacity    int    `json:"capacity"`
	OwnerUserID string `json:"owner_user_id"`
}
