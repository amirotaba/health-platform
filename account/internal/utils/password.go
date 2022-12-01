package utils

import (
	"math/rand"
	"strings"
)

var (
	// lowerCharSet   = "abcdedfghijklmnopqrst"
	lowerCharSet           = "abcdedfghijkmnopqrst"
	lowerVowelsCharSet     = "aeiou"
	lowerConsonantsCharSet = "bcddfghkmnpqrst"
	upperVowelsCharSet     = "AEU"
	upperConsonantsCharSet = "BCDFGHJKLMNPQRSTVWXYZ"
	upperCharSet           = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
	specialCharSet         = "!@#$%&*"
	numberSet              = "0123456789"
	allCharSet             = lowerCharSet + upperCharSet + specialCharSet + numberSet
	otherCharSet           = lowerCharSet
)

func GeneratePassword() string {
	// minSpecialChar := 1
	// minNum := 1
	// minUpperCase := 1
	// minLowerCase := 3
	// passwordLength := 8
	template := "baB0ba0X"
	// password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)
	password := generatePassword2(template)

	return password
}

func generatePassword2(template string) string {
	var password strings.Builder

	// template := "baB0ba0A0"
	for _, v := range template {
		switch v {
		case 'a':
			random := rand.Intn(len(lowerVowelsCharSet))
			password.WriteString(string(lowerVowelsCharSet[random]))
		case 'b':
			random := rand.Intn(len(lowerConsonantsCharSet))
			password.WriteString(string(lowerConsonantsCharSet[random]))
		case 'x':
			random := rand.Intn(len(lowerCharSet))
			password.WriteString(string(lowerCharSet[random]))
		case 'A':
			random := rand.Intn(len(upperVowelsCharSet))
			password.WriteString(string(upperVowelsCharSet[random]))
		case 'B':
			random := rand.Intn(len(upperConsonantsCharSet))
			password.WriteString(string(upperConsonantsCharSet[random]))
		case 'X':
			random := rand.Intn(len(upperCharSet))
			password.WriteString(string(upperCharSet[random]))
		case '0':
			random := rand.Intn(len(numberSet))
			password.WriteString(string(numberSet[random]))
		case '!':
			random := rand.Intn(len(specialCharSet))
			password.WriteString(string(specialCharSet[random]))
		}
	}
	// //Set special character
	// for i := 0; i < minSpecialChar; i++ {
	// 	random := rand.Intn(len(specialCharSet))
	// 	password.WriteString(string(specialCharSet[random]))
	// }

	// //Set numeric
	// for i := 0; i < minNum; i++ {
	// 	random := rand.Intn(len(numberSet))
	// 	password.WriteString(string(numberSet[random]))
	// }

	// //Set uppercase
	// for i := 0; i < minUpperCase; i++ {
	// 	random := rand.Intn(len(upperCharSet))
	// 	password.WriteString(string(upperCharSet[random]))
	// }

	//Set lowercase
	// for i := 0; i < minLowerCase; i++ {
	// 	random := rand.Intn(len(lowerCharSet))
	// 	password.WriteString(string(lowerCharSet[random]))
	// }

	// remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	// for i := 0; i < remainingLength; i++ {
	// 	random := rand.Intn(len(otherCharSet))
	// 	password.WriteString(string(otherCharSet[random]))
	// }
	inRune := []rune(password.String())
	// rand.Shuffle(len(inRune), func(i, j int) {
	// 	inRune[i], inRune[j] = inRune[j], inRune[i]
	// })
	return string(inRune)
}
