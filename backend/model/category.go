package model

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	IdleDevice  int    `json:"idle_device"`
	UsingDevice int    `json:"using_device"`
}

type ResponseGetCategories struct {
	Categories []Category `json:"categories"`
}

type ResponseGetCategory Category

type RequestCreateCategory Category

type ResponseCreateCategory struct {
	Message string `json:"message"`
}

type ResponseDeleteCategory struct {
	Message string `json:"message"`
}
