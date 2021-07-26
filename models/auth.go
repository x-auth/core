package models

type Profile struct {
	Name   string   `json:"first_name"`
	Email  string   `json:"email"`
	Groups []string `json:"group"`
}
