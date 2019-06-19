package models

type Zoo struct {
	Base
	Name     string `json:"name"`
	Location string `json:"location"`
}
