package model

import "time"

type History struct {
	User       string     `json:"user"`
	LastBorrow *time.Time `json:"lastBorrow"`
	LastReturn *time.Time `json:"lastReturn"`
}

type Histories []History

type ResponseGetHistory struct {
	Histories Histories
}
