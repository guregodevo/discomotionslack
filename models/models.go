package models

type PlaySearchItem struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type PlayVideo struct {
	VideoId string `json:"video_xid"`
	Now     bool   `json:"now"`
}
