package app

import (
	"database/sql"
	"errors"
	"fmt"
	"mastodon-archiver/internal/dao"
	"mastodon-archiver/internal/db"
	"mastodon-archiver/internal/mastodon"
	"mastodon-archiver/internal/secret"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/perimeterx/marshmallow"
)

type App struct {
	PostBlobDao    *dao.PostBlobDao
	RunHistoryDao  *dao.RunHistoryDao
	MastodonClient *mastodon.MastodonClient
}

func InitApp(secret *secret.Secret) (*App, error) {
	dbx, err := db.Connect(secret.DbHost, secret.DbPort, secret.DbUser, secret.DbPassword, secret.DbName)
	if err != nil {
		return nil, err
	}
	postBlobDao := dao.NewPostBlobDao(dbx)
	runHistoryDao := dao.NewRunHistoryDao(dbx)

	client := mastodon.NewMastodonClient("https://social.mikutter.hachune.net", secret.MastodonAccessToken)

	return &App{
		PostBlobDao:    postBlobDao,
		RunHistoryDao:  runHistoryDao,
		MastodonClient: client,
	}, nil
}

func (app *App) RunRecurring() error {
	startedAt := time.Now()
	runId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	maxSeenStatusId, err := app.RunHistoryDao.QueryMaxSeenStatusId()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			maxSeenStatusId = "0"
		} else {
			return err
		}
	}

	fmt.Printf("Start from ID=%s\n", maxSeenStatusId)
	postStrings, err := app.MastodonClient.GetHomeRaw(maxSeenStatusId)
	if err != nil {
		return err
	}
	fmt.Printf("Got %d posts\n", len(postStrings))

	postBlobs := make([]dao.PostBlob, len(postStrings))
	for _, postString := range postStrings {
		var post mastodon.Status
		if _, err := marshmallow.Unmarshal([]byte(postString), &post); err != nil {
			println(postString)
			return err
		}
		if post.Id > maxSeenStatusId {
			maxSeenStatusId = post.Id
		}

		postBlob := dao.PostBlob{StatusId: post.Id, JsonBody: postString}
		if err = app.PostBlobDao.Insert(&postBlob); err != nil {
			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
				fmt.Printf("Skip duplicate post %s\n", post.Id)
				continue
			} else {
				return err
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
	app.RunHistoryDao.Insert(&runHistory)
	fmt.Printf("%+v\n", runHistory)

	return nil
}

func (app *App) RunRecurringForever() error {
	for {
		if err := app.RunRecurring(); err != nil {
			return err
		}
		time.Sleep(60 * time.Second)
	}
}
