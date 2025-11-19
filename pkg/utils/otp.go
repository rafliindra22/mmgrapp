package utils

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() (string, string, error) {
	b := make([]byte, 3) // 3 byte = 24 bit
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}

	// convert random bytes to 6 digit number
	num := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	otp := fmt.Sprintf("%06d", num) // selalu 6 digit, leading zero included

	// hash OTP untuk simpan
	hashedOTPBytes, _ := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	hashedOTP := string(hashedOTPBytes)

	return otp, hashedOTP, nil
}
