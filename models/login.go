package models

type Login struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
