package model

type RequestQrcodeBorrow struct {
	User string `json:"user"`
}

type RequestQrcodeReturn struct {
	User string `json:"user"`
}
