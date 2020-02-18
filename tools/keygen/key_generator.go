package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
)

func main() {
	length := flag.Int("len", 32, "")
	flag.Parse()

	b := make([]byte, *length)
	_, err := rand.Read(b)
	if err == nil {
		result := base64.StdEncoding.EncodeToString(b)
		if len(result) > *length {
			result = result[:*length]
		}
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
}
