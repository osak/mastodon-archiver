package main

import (
	"os"
	"mastodon-archiver/internal/db"
	"mastodon-archiver/internal/dao"
	"mastodon-archiver/internal/mastodon"
	"mastodon-archiver/internal/secret"
	"github.com/google/uuid"
)

func main() {
	secret, err := secret.LoadSecret(os.Args[1])
	if err != nil {
		panic(err)
	}

	dbx, err := db.Connect(secret.DbHost, secret.DbPort, secret.DbUser, secret.DbPassword, secret.DbName)
	if err != nil {
	    panic(err)
	}
	postBlobDao := dao.NewPostBlobDao(dbx)

	client := mastodon.NewMastodonClient("https://social.mikutter.hachune.net", secret.MastodonAccessToken)
	postStrings, err := client.GetHomeRaw()
	if err != nil {
	    panic(err)
	}

	postBlobs := make([]dao.PostBlob, len(postStrings))
	for _, postString := range postStrings {
	    uuid, err := uuid.NewV7()
	    if err != nil {
		panic(err)
	    }
	    postBlob := dao.PostBlob{Id: uuid.String(), JsonBody: postString}
	    postBlobDao.Insert(&postBlob)
	    postBlobs = append(postBlobs, postBlob)
	}
}
