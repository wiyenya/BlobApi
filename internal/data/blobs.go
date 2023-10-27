package data

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

type Blob struct {
	Index  int64          `db:"index"`
	UserId *int32         `db:"user_id" json:"user_id"`
	Data   types.JSONText `db:"data" json:"data"`
}

func (b Blob) MarshalJSON() ([]byte, error) {
	type Alias Blob
	return json.Marshal(&struct {
		*Alias
		UserId *int32         `json:"user_id,omitempty"`
		Data   types.JSONText `json:"data"`
	}{
		Alias:  (*Alias)(&b),
		UserId: b.UserId,
		Data:   b.Data,
	})
}

type Blobs interface {
	Insert(userID int, data map[string]interface{}) (int, error)
	Get(id int) (*Blob, error)
	GetBlobList() ([]*Blob, error)
	Delete(id int) error
}
