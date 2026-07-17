package model

type ErrorDetail struct {
	HttpStatus int    `json:"-"`
	Detail     string `json:"message"`
}
