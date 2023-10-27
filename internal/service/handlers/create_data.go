package handlers

import (
	data "BlobApi/internal/data"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx/types"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"

	// "gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/horizon-connector"
)

type DataCreate struct {
	Tx *xdrbuild.Transaction
	//InviteQ        data.Invite
	Horizon *horizon.Connector
	//SystemSettings data.SystemSettings
}

func (c DataCreate) CreateDataHandler(w http.ResponseWriter, r *http.Request) {

	tx := c.Tx

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

	tx = tx.Op(&createDataOp)

	// Signing

	// Get key
	SECRET_KEY := "SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4"
	kp, err := keypair.Parse(SECRET_KEY)
	if err != nil {
		http.Error(w, "Failed to parse secret key", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, tx)
	if err1 != nil {
		http.Error(w, "Failed to marshal transaction", http.StatusInternalServerError)
		return
	}
	txBytes := buf.Bytes()

	// Sign tx
	signedTransaction, err := kp.Sign(txBytes)
	if err != nil {
		http.Error(w, "Failed to sign transaction", http.StatusInternalServerError)
		return
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedTransaction)

	// Send tx to Horizon
	endpoint := "https://docs.tokend.io/horizon#operation/submitTransaction"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	if err != nil {
		http.Error(w, "Failed to submit transaction", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Responce
	response := map[string]interface{}{
		"transaction": encodedSignedTransaction,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
