package sms

import (
	"errors"
	"fmt"
	"net/url"

	"git.paygear.ir/giftino/account/internal/utils"
)

func SendSms(phoneNumber string, content string) error {
	from := "%2B9830000600612222"
	content = url.QueryEscape(content)
	urlStr := fmt.Sprintf("https://sms.magfa.com/magfaHttpService?service=Enqueue&domain=mashhad_ct&accountname=mancard&password=NYKZWAK4DELwEcgJ&from=%s&to=%s&message=%s", from, phoneNumber, content)
	method := "GET"

	accountGetLogResp, statusCode, err := utils.SMSRPCCall(
		urlStr,
		"",
		method,
		nil)

	if err != nil {
		return err
	}
	if statusCode != 200 {
		return errors.New(string(accountGetLogResp))
	}
	return nil
}

func GetMagfaBalance() (string, error) {
	res, status, err := utils.RPCCall("https://sms.magfa.com/api/http/sms/v1?service=getcredit&accountname=mancard&password=NYKZWAK4DELwEcgJ&domain=mashhad_ct", "", "GET", nil)
	if err != nil {
		return "", err
	}
	if status == 200 {
		resStr := string(res)
		return resStr, nil
	}
	return "", errors.New("error in getting balance")
}
