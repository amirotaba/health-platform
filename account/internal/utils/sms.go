package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type Payload struct {
	Token    string `json:"Token"`
	Template string `json:"Template"`
	Mobile   string `json:"Mobile"`
	Message  string `json:"Message"`
}

func SMS(request entity.SendSmsRequest) error {
	data := Payload{
		Token:    domain.SMSToken,
		Template: domain.SMSTemplate,
		Mobile:   request.Mobile,
		Message:  request.Message,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://sanasp.ir/api/v1/SendOTPCode", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
