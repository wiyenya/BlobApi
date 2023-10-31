package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/strkey"
	"gitlab.com/tokend/horizon-connector"
	"net/http"

	dataPkg "BlobApi/internal/data"

	"github.com/jmoiron/sqlx/types"

	"BlobApi/internal/service"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

// Define the structure to process the response
type ServerResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type HorizonModel struct {
	log     *logan.Entry
	horizon *horizon.Connector
}

var s = service.NewImportant(
	"SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4",
	"http://localhost:8000/_/api/",
	"http://localhost:8000/_/api/v3/data/",
	"http://localhost:8000/_/api/v3/data/",
	"TokenD Developer Network",
	601200,
)

func keyP(str string) *keypair.Full {
	//for Singing
	raw, err := strkey.Decode(strkey.VersionByteSeed, str)
	if err != nil {
		panic(err)
	}

	var seed [32]byte
	if len(raw) != 32 {
		panic("decoded seed is not 32 bytes long")
	}
	copy(seed[:], raw)

	kp, err := keypair.FromRawSeed(seed)
	return kp
}

func (q *HorizonModel) Insert(userID int32, data types.JSONText) (int, error) {

	blob := dataPkg.Blob{
		Index:  1,
		UserId: &userID,
		Data:   data,
	}

	createData := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	address, err := keypair.Parse(s.Seed)

	if err != nil {
		panic(err)
	}
	builder := xdrbuild.NewBuilder(s.NetworkPassphrase, s.TxExpirationPeriod)
	tx := builder.Transaction(address)

	txEnvelope := tx.Op(createData)
	if err != nil {
		panic(err)
	}

	//Singing
	signedEnvelope := txEnvelope.Sign(keyP(s.Seed))

	encodedSignedTransaction, err := xdr.MarshalBase64(signedEnvelope)

	resp, err := http.Post(s.EndpointForPost, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
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
		return -1, fmt.Errorf("server returned non-success status")
	}

	return response.ID, nil

}

func (q *HorizonModel) Get(id int) (*dataPkg.Blob, error) {

	str := s.EndpointForGet + "%d"
	endpoint := fmt.Sprintf(str, id)

	response, err := q.horizon.Client().Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var blob dataPkg.Blob
	if err := json.Unmarshal(response, &blob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &blob, nil

}

func (q *HorizonModel) GetBlobList() ([]*dataPkg.Blob, error) {

	response, err := q.horizon.Client().Get(s.EndpointForGetList)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var blobs []*dataPkg.Blob
	if err := json.Unmarshal(response, &blobs); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	var blobs2 []*dataPkg.Blob
	for _, blob := range blobs {
		blobs2 = append(blobs2, blob)
	}

	return blobs2, nil
}

func (q *HorizonModel) Delete(id int) error {
	removeData := xdrbuild.RemoveData{
		ID: uint64(id),
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

	//Singing
	signedEnvelope := txEnvelope.Sign(keyP(s.Seed))

	encodedSignedTransaction, err := xdr.MarshalBase64(signedEnvelope)

	endpoint := "http://localhost:8000/_/api/"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	// if err != nil {
	// 	return "", err
	// }
	defer resp.Body.Close() //guaranteed function fulfillment

	return nil

}
