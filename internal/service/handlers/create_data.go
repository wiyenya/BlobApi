package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"

	// "gitlab.com/tokend/horizon-connector"
	data "BlobApi/internal/data"

	"gitlab.com/tokend/horizon-connector"
)

type DataCreate struct {
	Tx *xdrbuild.Transaction
	//InviteQ        data.Invite
	Horizon *horizon.Connector
	//SystemSettings data.SystemSettings
}

func (c DataCreate) CreateDataHandler(w http.ResponseWriter, r *http.Request) {

	tx := data.Transaction()
	// Signing
	SECRET_KEY := "SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4"

	signedTransaction, err := data.Signing(SECRET_KEY, tx)
	if err != nil {
		Log(r).WithError(err).Error("Failed to sign transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	encodedSignedTransaction, err2 := data.Encode(signedTransaction)
	if err2 != nil {
		Log(r).WithError(err).Error("Failed to Marshal")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Send tx to Horizon
	endpoint := "http://localhost:8000/_/api/"
	resp, err3 := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	if err3 != nil {
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
