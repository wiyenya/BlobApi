package data

type Blob struct {
	Index  int64  `db:"index"`
	UserId *int32 `db:"user_id" json:"user_id"`
	Data   []byte `db:"data" json:"data"`
}

type Blobs interface {
	Insert(userID int, data map[string]interface{}) (int, error)
	Get(id int) (*Blob, error)
	GetBlobList() ([]*Blob, error)
	Delete(id int) error
}
