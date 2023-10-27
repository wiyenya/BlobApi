package handlers

import (
	data "BlobApi/internal/data"
	"bytes"
	"net/http"

	"github.com/jmoiron/sqlx/types"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

func blobTr(tx *xdrbuild.Transaction) (string, error) {

	userID := int32(42)
	blob := data.Blob{
		Index:  1,
		UserId: &userID,
		Data:   types.JSONText(`{"name":"John", "age":30}`),
	}

	createDataOp := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	tx = tx.Op(createDataOp)

	// Signing

	// Get the key
	SECRET_KEY := "SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4"
	kp, err := keypair.Parse(SECRET_KEY)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, tx)
	if err1 != nil {
		panic(err1)
	}
	txBytes := buf.Bytes()

	// Sign the transaction
	signedTransaction, err := kp.Sign(txBytes)
	if err != nil {
		panic(err)
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedTransaction)

	endpoint := "https://docs.tokend.io/horizon#operation/submitTransaction"

	// http.Post ожидает третьим аргументом значение типа io.Reader, а у вас переменная encodedSignedTransaction типа string
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return encodedSignedTransaction, nil
}
