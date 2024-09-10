package models


type Url struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type MultipleFileUploadResponse struct {
	Url []*Url `json:"urls"`
}
