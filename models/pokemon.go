package models

type Pokemon struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Level int    `json:"level"`
}
