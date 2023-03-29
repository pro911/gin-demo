package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func main() {

	timestamp := time.Now().Unix()
	timestamp = 1679662370
	secret := "PxCjg9ll0NHUwmhsRCWMee"
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	sign, _ := GenSign(secret, timestamp)
	fmt.Println(stringToSign, timestamp, secret, sign)
}
