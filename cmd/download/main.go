package main

import (
	"os"
	"mastodon-archiver/internal/mastodon"
	"mastodon-archiver/internal/secret"
)

func main() {
	secret, err := secret.LoadSecret(os.Args[1])
	if err != nil {
		panic(err)
	}

	client := mastodon.NewMastodonClient("https://social.mikutter.hachune.net", secret.MastodonAccessToken)
	postStrings, err := client.GetHomeRaw()
	if err != nil {
	    panic(err)
	}

	for _, postString := range postStrings {
	    println(postString)
	}
}
