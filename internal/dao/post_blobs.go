import (
    "github.com/jmoiron/sqlx"
)

type PostBlob struct {
    Id string 
    JsonBody string
}

type PostBlobDao struct {}

func (dao *PostBlobDao) Insert(db *sqlx.DB, postBlob *PostBlob) error {
	_, err := db.Exec("INSERT INTO post_blobs (id, json_body) VALUES (?, ?)", postBlob.Id, postBlob.JsonBody)
	return err
}
