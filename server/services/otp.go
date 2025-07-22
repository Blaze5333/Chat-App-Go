package services

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
