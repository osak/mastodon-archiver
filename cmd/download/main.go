package main

import (
	"database/sql"
	"errors"
	"fmt"
	"mastodon-archiver/internal/dao"
	"mastodon-archiver/internal/db"
	"mastodon-archiver/internal/mastodon"
	"mastodon-archiver/internal/secret"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/perimeterx/marshmallow"
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
	runHistoryDao := dao.NewRunHistoryDao(dbx)

	startedAt := time.Now()
	runId, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	maxSeenStatusId, err := runHistoryDao.QueryMaxSeenStatusId()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			maxSeenStatusId = "0"
		} else {
			panic(err)
		}
	}

	fmt.Printf("Start from ID=%s\n", maxSeenStatusId)
	client := mastodon.NewMastodonClient("https://social.mikutter.hachune.net", secret.MastodonAccessToken)
	postStrings, err := client.GetHomeRaw(maxSeenStatusId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got %d posts\n", len(postStrings))

	postBlobs := make([]dao.PostBlob, len(postStrings))
	for _, postString := range postStrings {
		var post mastodon.Status
		if _, err := marshmallow.Unmarshal([]byte(postString), &post); err != nil {
			println(postString)
			panic(err)
		}
		if post.Id > maxSeenStatusId {
			maxSeenStatusId = post.Id
		}

		postBlob := dao.PostBlob{StatusId: post.Id, JsonBody: postString}
		if err = postBlobDao.Insert(&postBlob); err != nil {
			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
				fmt.Printf("Skip duplicate post %s\n", post.Id)
				continue
			} else {
				panic(err)
			}
		}
		postBlobs = append(postBlobs, postBlob)
	}
	finishedAt := time.Now()
	runHistory := dao.RunHistory{
		Id: runId.String(),
		RunType: dao.RunTypeRecurring,
		StartedAt: startedAt,
		FinishedAt: finishedAt,
		MaxSeenStatusId: maxSeenStatusId,
	}
	runHistoryDao.Insert(&runHistory)
	fmt.Printf("%+v\n", runHistory)
}
