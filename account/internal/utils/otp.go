package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func Otp() string {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return strconv.Itoa(rand.Intn(max-min) + min)

}
