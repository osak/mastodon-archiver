package main

import (
	"mastodon-archiver/internal/app"
	"mastodon-archiver/internal/secret"
	"os"
)

func main() {
	secret, err := secret.LoadSecret(os.Args[1])
	if err != nil {
		panic(err)
	}
	app, err := app.InitApp(secret)
	if err != nil {
		panic(err)
	}

	if err = app.RunRecurringForever(); err != nil {
		panic(err)
	}
}
