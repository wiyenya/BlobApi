package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	// data "BlobApi/internal/data"
	// "BlobApi/internal/createdata"

	// "github.com/jmoiron/sqlx/types"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	// "gitlab.com/tokend/go/xdrbuild"
	// "gitlab.com/tokend/horizon-connector"
)

func createDataHandler(w http.ResponseWriter, r *http.Request) {

	dataCreate := DataCreate{
		// Initialize with the appropriate values
		Tx:      initializeTransaction(),
		Horizon: initializeHorizonConnector(),
	}

	tx := DataCreate(r).CreateData(r.Context())

	// Signing

	// Получаем ключ
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

	// Подписываем транзакцию
	signedTransaction, err := kp.Sign(txBytes)
	if err != nil {
		http.Error(w, "Failed to sign transaction", http.StatusInternalServerError)
		return
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedTransaction)

	// Отправляем транзакцию на сервер Horizon
	endpoint := "https://docs.tokend.io/horizon#operation/submitTransaction"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	if err != nil {
		http.Error(w, "Failed to submit transaction", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Возвращаем ответ клиенту, например, в виде JSON
	response := map[string]interface{}{
		"transaction": encodedSignedTransaction,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
