package data

type Blob struct {
	ID     int    `json:"id"`
	UserID *int32 `json:"user_id"`
	Data   []byte `json:"data"`
}
