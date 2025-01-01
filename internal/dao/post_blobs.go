package dao

import (
    "github.com/jmoiron/sqlx"
)

type PostBlob struct {
    Id string 
    JsonBody string
}

type PostBlobDao struct {
    dbx *sqlx.DB
}

func NewPostBlobDao(dbx *sqlx.DB) *PostBlobDao {
	return &PostBlobDao{dbx: dbx}
}

func (dao *PostBlobDao) Insert(postBlob *PostBlob) error {
	_, err := dao.dbx.Exec("INSERT INTO post_blobs (id, json_body) VALUES (?, ?)", postBlob.Id, postBlob.JsonBody)
	return err
}
