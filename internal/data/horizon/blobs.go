package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	regources "gitlab.com/tokend/regources/generated"

	dataPkg "BlobApi/internal/data"

	"github.com/jmoiron/sqlx/types"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
)

// Define the structure to process the response
type ServerResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type HorizonModel struct {
	log     *logan.Entry
	horizon *horizon.Connector
	kp      keypair.Full
	builder *xdrbuild.Builder
}

func NewHorizonModel(log *logan.Entry, domain string, seed string) *HorizonModel {
	kp := keypair.MustParseSeed(seed)
	horizonUrl, err := url.Parse(domain)
	if err != nil {
		panic(err)
	}

	horizonClient := horizon.NewConnector(horizonUrl).
		WithSigner(kp)

	txBuilder, err := horizonClient.TXBuilder()
	if err != nil {
		panic(err)
	}

	return &HorizonModel{
		log:     log,
		horizon: horizonClient,
		kp:      kp,
		builder: txBuilder,
	}
}

var ErrBlobNotFound = errors.New("blob not found")

func (q *HorizonModel) Insert(userID int32, data types.JSONText) (int, error) {

	blob := dataPkg.Blob{
		UserId: &userID,
		Data:   data,
	}

	createData := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	tx, err := q.builder.Transaction(q.kp).
		Op(createData).
		Sign(q.kp).
		Marshal()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build tx")
	}

	resp := q.horizon.Submitter().Submit(context.TODO(), tx)

	var txResponse xdr.TransactionResult
	err = xdr.SafeUnmarshalBase64(resp.ResultXDR, &txResponse)
	if err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal tx response")
	}

	// 	err = xdr.SafeUnmarshalBase64(result.ResultXDR, txResp)

	var resultID int

	for _, err1 := range txResponse.Result.MustResults() {
		a, err2 := err1.MustTr().GetCreateDataResult()
		if !err2 {
			return 0, errors.Wrap(err, "failed to get create data result")
		}

		resultID = int(a.Success.DataId)
		break
	}

	return resultID, nil

}

func (q *HorizonModel) Get(id int) (*dataPkg.Blob, error) {

	resp, err := q.horizon.Client().Get(fmt.Sprintf("/v3/data/%d", id))

	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if resp == nil && err == nil {
		return nil, ErrBlobNotFound
	}

	var parsedResponse regources.DataResponse
	if err := json.Unmarshal(resp, &parsedResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	var blob dataPkg.Blob
	if err := json.Unmarshal(parsedResponse.Data.Attributes.Value, &blob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	blob.Index = int64(id)

	return &blob, nil

}

func (q *HorizonModel) GetBlobList() ([]*dataPkg.Blob, error) {

	response, err := q.horizon.Client().Get("/v3/data")

	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var parsedResponse regources.DataListResponse
	if err := json.Unmarshal(response, &parsedResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	result := make([]*dataPkg.Blob, len(parsedResponse.Data))
	for i, data := range parsedResponse.Data {
		var blob dataPkg.Blob
		if err := json.Unmarshal(data.Attributes.Value, &blob); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		result[i] = &blob
	}

	return result, nil
}

func (q *HorizonModel) Delete(id int) error {

	//try to find blob my id
	resp, err := q.horizon.Client().Get(fmt.Sprintf("/v3/data/%d", id))
	if len(resp) == 0 {
		customErr := fmt.Errorf("Blob not found for ID %d", id)
		return errors.Wrap(customErr, "request failed")
	}

	logan.New().Debug("test")
	removeData := xdrbuild.RemoveData{
		ID: uint64(id),
	}

	tx, err := q.builder.Transaction(q.kp).
		Op(removeData).
		Sign(q.kp).
		Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to build tx")
	}

	q.horizon.Submitter().Submit(context.TODO(), tx)

	return nil

}
