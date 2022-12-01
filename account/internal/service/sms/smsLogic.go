package sms

import (
	"fmt"
)

const ServiceContentName = ""

func SendSmsAsync(calledName string, transactionId *uint, MobilePhones []string, content string, orderId *uint) {
	if orderId != nil {
		content = fmt.Sprintf(ServiceContentName+"WorkerName: %s\nOrderId : %d\nError: %s", calledName, *orderId, content)
	} else if transactionId != nil {
		content = fmt.Sprintf(ServiceContentName+"WorkerName: %s\nTransactionId : %d\nError: %s", calledName, *transactionId, content)
	}
	go func(mobilePhones []string) {
		//sendSms
		for _, v := range mobilePhones {
			err := SendSms(v, content)
			if err != nil {
				err = SendSms(v, content)
				if err != nil {
					err = SendSms(v, content)
					if err != nil {
						fmt.Printf("Failed to send sms to : %s with content : %s", mobilePhones, content)
						continue
					}
				}
			}
		}
	}(MobilePhones)
}
