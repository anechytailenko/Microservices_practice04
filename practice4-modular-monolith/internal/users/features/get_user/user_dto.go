package get_user

type UserDTO struct {
	ID          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	Meetups     []string `json:"meetups"`
}
