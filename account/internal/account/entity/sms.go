package entity

type SendSmsRequest struct {
	Mobile  string `json:"Mobile"`
	Message string `json:"Message"`
}

type Sanasp struct {
	Token    string `json:"Token"`
	Template string `json:"Template"`
}
