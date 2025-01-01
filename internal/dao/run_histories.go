package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type RunType int

const (
	RunTypeRecurring = 1
	RunTypeOneOff    = 2
)

type RunHistory struct {
	Id              string
	RunType         RunType
	StartedAt       time.Time
	FinishedAt      time.Time
	MaxSeenStatusId string
}

type RunHistoryDao struct {
	dbx *sqlx.DB
}

func NewRunHistoryDao(dbx *sqlx.DB) *RunHistoryDao {
	return &RunHistoryDao{dbx: dbx}
}

func (dao *RunHistoryDao) Insert(runHistory *RunHistory) error {
	_, err := dao.dbx.Exec(
		"INSERT INTO run_histories (id, run_type, started_at, finished_at, max_seen_status_id) VALUES (?, ?, ?, ?, ?)",
		runHistory.Id, runHistory.RunType, runHistory.StartedAt, runHistory.FinishedAt, runHistory.MaxSeenStatusId,
	)
	return err
}

func (dao *RunHistoryDao) QueryMaxSeenStatusId() (string, error) {
	var maxSeenStatusId string
	err := dao.dbx.Get(&maxSeenStatusId, "SELECT max_seen_status_id FROM run_histories WHERE run_type = ? ORDER BY started_at DESC LIMIT 1", RunTypeRecurring)
	return maxSeenStatusId, err
}
