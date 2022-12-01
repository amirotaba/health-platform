package entity

type Tag struct {
	Exist        bool   `json:"Exist"`
	Status       bool   `json:"Status"`
	Name         string `json:"Name"`
	ID           int64  `json:"ID"`
	CategoryName string `json:"CategoryName"`
	CategoryID   int64  `json:"CategoryID"`
}
