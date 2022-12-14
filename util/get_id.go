package util

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"math/rand"
	"strings"
	"time"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//GetID represents get unique id for transaction
func GetID() string {
	dest, _ := hex.DecodeString(fmt.Sprintf("%d", nowAsUnixSecond()))
	var id strings.Builder
	encode := base64.StdEncoding.EncodeToString(dest)
	rand.Seed(time.Now().UnixNano())
	id.WriteString(encode)
	id.WriteString(randString(4))
	return strings.Replace(id.String(), "=", randString(1), 1)
}

func nowAsUnixSecond() int64 {
	return time.Now().UnixNano() / 1e9
}
