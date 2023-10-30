package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	dataPkg "BlobApi/internal/data"

	"github.com/jmoiron/sqlx/types"

	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

// Define the structure to process the response
type ServerResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func Insert(userID int32, data types.JSONText) (string, error) {

	blob := dataPkg.Blob{
		Index:  1,
		UserId: &userID,
		Data:   data,
	}

	createData := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	address, err := keypair.Parse("SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")
	if err != nil {
		panic(err)
	}
	builder := xdrbuild.NewBuilder("<NETWORK_PASSPHRASE>", 300)
	tx := builder.Transaction(address)

	txEnvelope := tx.Op(createData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, txEnvelope)
	if err1 != nil {
		panic(err1)
	}
	txBytes := buf.Bytes()

	signedEnvelope, err := address.Sign(txBytes)
	if err != nil {
		panic(err)
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedEnvelope)

	endpoint := "http://localhost:8000/_/api/"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	// if err != nil {
	// 	return "", err
	// }
	defer resp.Body.Close() //guaranteed function fulfillment

	// Response processing
	var response ServerResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	if response.Status != "success" {
		return "", fmt.Errorf("server returned non-success status")
	}

	return response.ID, nil

}

func Delete(userID int32) (string, error) {

	removeData := xdrbuild.RemoveData{
		ID: uint64(userID),
	}

	address, err := keypair.Parse("SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")
	if err != nil {
		panic(err)
	}
	builder := xdrbuild.NewBuilder("<NETWORK_PASSPHRASE>", 300)
	tx := builder.Transaction(address)

	txEnvelope := tx.Op(removeData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, txEnvelope)
	if err1 != nil {
		panic(err1)
	}
	txBytes := buf.Bytes()

	signedEnvelope, err := address.Sign(txBytes)
	if err != nil {
		panic(err)
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedEnvelope)

	endpoint := "http://localhost:8000/_/api/"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	// if err != nil {
	// 	return "", err
	// }
	defer resp.Body.Close() //guaranteed function fulfillment

	return encodedSignedTransaction, nil
}

func Update(userID int32, data types.JSONText) (string, error) {
	updateData := xdrbuild.UpdateData{
		ID:    uint64(userID),
		Value: data,
	}

	address, err := keypair.Parse("SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")
	if err != nil {
		panic(err)
	}
	builder := xdrbuild.NewBuilder("<NETWORK_PASSPHRASE>", 300)
	tx := builder.Transaction(address)

	txEnvelope := tx.Op(updateData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	_, err1 := xdr.Marshal(&buf, txEnvelope)
	if err1 != nil {
		panic(err1)
	}
	txBytes := buf.Bytes()

	signedEnvelope, err := address.Sign(txBytes)
	if err != nil {
		panic(err)
	}

	encodedSignedTransaction, err := xdr.MarshalBase64(signedEnvelope)

	endpoint := "http://localhost:8000/_/api/"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	// if err != nil {
	// 	return "", err
	// }
	defer resp.Body.Close() //guaranteed function fulfillment

	return encodedSignedTransaction, nil
}

const horizonBaseUrl = "https://horizon.tokend.io/v3"

// func Get(id int) (*dataPkg.Blob, error) {
// 	url := fmt.Sprintf("%s/data/%d", horizonBaseUrl, id)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("error fetching data from Horizon: %s", resp.Status)
// 	}

// 	var blob dataPkg.Blob
// 	err = json.NewDecoder(resp.Body).Decode(&blob)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &blob, nil
// }

// func GetBlobList() ([]*dataPkg.Blob, error) {
// 	url := fmt.Sprintf("%s/data", horizonBaseUrl)

// 	// Добавьте параметры, если необходимо, например:
// 	// url += "?filter[type]=YOUR_TYPE&filter[owner]=YOUR_OWNER"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("error fetching data from Horizon: %s", resp.Status)
// 	}

// 	var blobs []*dataPkg.Blob
// 	err = json.NewDecoder(resp.Body).Decode(&blobs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return blobs, nil
// }
