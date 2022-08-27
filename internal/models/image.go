package models

type ImageResizeOption struct {
	Width      int `json:"saveOptions"`
	Height     int `json:"height"`
	Percentage int `json:"percentage"`
}

type ImageRequest struct {
	FileName string `json:"fileName" validator:"required"`
	File     string `json:"file" validator:"required"`
}

type ImageResizeRequest struct {
	Images       []ImageRequest    `json:"images"`
	ResizeType   string            `json:"resizeType"`
	ResizeOption ImageResizeOption `json:"resizeOption"`
	SaveOptions  string            `json:"saveOptions"`
}
