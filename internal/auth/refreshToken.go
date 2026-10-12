package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToekn() string {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		log.Println(err)
	}
	token := hex.EncodeToString(data)
	return token
}
