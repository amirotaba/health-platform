package utils

func CheckPhone(phone string) (string, error) {
	//match, _ := regexp.MatchString("^(\\+98|0)?9\\d{9}$", phone)
	//if !match {
	//	return "", errors.New("phone number invalid")
	//}
	//
	//if strings.HasPrefix(phone, "+98") {
	//	return fmt.Sprintf("0%s", phone[3:]), nil
	//}

	return phone, nil
}
