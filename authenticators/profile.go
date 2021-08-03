package authenticators

type Profile struct {
	Id     string   `json:"id"`
	Name   string   `json:"first_name"`
	Email  string   `json:"email"`
	Groups []string `json:"group"`
}
