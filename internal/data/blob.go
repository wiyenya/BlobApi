package data

type Blob struct {
	ID     int    `json:"id"`
	UserID *int32 `json:"user_id"`
	Data   string `json:"data"`
}

type BlobStorer interface {
	Insert(userID int, data string) (int, error)
	Get(id int) (*Blob, error)
	GetBlobList() ([]*Blob, error)
	Delete(id int) error
}
