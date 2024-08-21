package helper

import (
	"crypto/aes"
	"log"
	"os"
)

func EncryptRefreshToken(refresh string) string {
	aesKey := os.Getenv("AES_KEY")
	cipher, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	dst := make([]byte, len(refresh))
	cipher.Encrypt(dst, []byte(refresh))
	return string(dst)
}

func DecryptRefreshToken(encrypted string) string {
	aesKey := os.Getenv("AES_KEY")
	cipher, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	dst := make([]byte, len(encrypted))
	cipher.Decrypt(dst, []byte(encrypted))
	return string(dst)
}

func GetTime(refresh string) string {
	ans := make([]byte, 0)
	for _, e := range refresh {
		if e == '.' {
			break
		}
		ans = append(ans, byte(e))
	}
	return string(ans)
}
