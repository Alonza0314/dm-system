package model

type Device struct {
	Id       string `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name" bind:"required"`
	Status   string `json:"status"`
	User     string `json:"user"`
	Owner    string `json:"owner"`
	Note     string `json:"note"`
}
