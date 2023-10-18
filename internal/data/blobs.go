package data

type Blob struct {
	Index   int    `json:"id"`
	User_id *int32 `json:"user_id"`
	Data    []byte `json:"data"`
}

type Blob2 struct {
	Index   int                    `json:"id"`
	User_id *int32                 `json:"user_id"`
	Data    map[string]interface{} `json:"data"`
}

type Blobs interface {
	Insert(userID int, data map[string]interface{}) (int, error)
	Get(id int) (*Blob2, error)
	GetBlobList() ([]*Blob2, error)
	Delete(id int) error
}
