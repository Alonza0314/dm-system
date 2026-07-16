package model

type Device struct {
	Id       int    `json:"id"`
	Category string `json:"category" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Status   string `json:"status"`
	User     string `json:"user"`
	Owner    string `json:"owner"`
	Note     string `json:"note"`
}

type DeviceShort struct {
	Id     int    `json:"id"`
	Name   string `json:"name" bind:"required"`
	Status string `json:"status"`
	User   string `json:"user"`
}

type ResponseGetDevices struct {
	Devices []DeviceShort
}

type ResponseGetDevice Device

type RequestCreateDevice Device

type ResponseCreateDevice struct {
	Message string `json:"message"`
}

type ResponseDeleteDevice struct {
	Message string `json:"message"`
}
