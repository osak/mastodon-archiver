package main

import (
	"os"
)

func main() {
	secret, err := LoadSecret(os.Args[1])
	if err != nil {
		panic(err)
	}

	client := NewMastodonClient("https://social.mikutter.hachune.net", secret.MastodonAccessToken)
	client.GetHome()
}
