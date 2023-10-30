package data

import (
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/tokend/go/xdrbuild"
)

func Transaction() xdrbuild.Transaction {

	tx := xdrbuild.NewBuilder("test", int64(1234)).Transaction("SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")

	userID := int32(42)
	blob := Blob{
		Index:  1,
		UserId: &userID,
		Data:   types.JSONText(`{"name":"John", "age":30}`),
	}

	createDataOp := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	tx = tx.Op(&createDataOp)

	return *tx
}
