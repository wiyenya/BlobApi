package data

import (
	"bytes"
	"fmt"

	"github.com/jmoiron/sqlx/types"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

func Transaction() xdrbuild.Transaction {

	tx := xdrbuild.NewBuilder("test", int64(1234)).Transaction()

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

func Signing(key string, tx xdrbuild.Transaction) ([]byte, error) {

	kp, err := keypair.Parse(key)
	if err != nil {
		fmt.Print(err)
		fmt.Print("Failed to parse secret key")
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, tx)

	fmt.Print(err1)
	if err1 != nil {
		fmt.Print(err1)
		fmt.Print("Failed to marshal transaction")
	}
	txBytes := buf.Bytes()

	// Sign tx
	bs := make([]byte, 0)
	signedTransaction, err := kp.Sign(txBytes)
	if err != nil {
		return bs, err
	}

	return signedTransaction, nil
}

func Encode(signedTransaction []byte) (string, error) {

	encodedSignedTransaction, err := xdr.MarshalBase64(signedTransaction)
	if err != nil {
		return encodedSignedTransaction, err
	}
	return encodedSignedTransaction, nil

}
