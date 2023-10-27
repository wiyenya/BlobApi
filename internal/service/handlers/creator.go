package handlers

import (
	"context"
	// "errors"
	// "fmt"

	data "BlobApi/internal/data"

	//validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx/types"
	// "gitlab.com/tokend/go/keypair"
	// "gitlab.com/tokend/go/xdr"
	// validation "github.com/go-ozzo/ozzo-validation"
	// "gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

// type DataCreate struct {
// 	Tx *xdrbuild.Transaction
// 	//InviteQ        data.Invite
// 	Horizon *horizon.Connector
// 	//SystemSettings data.SystemSettings
// }

func New(tx *xdrbuild.Transaction, horizon *horizon.Connector) DataCreate {
	return DataCreate{
		Tx:      tx,
		Horizon: horizon,
	}
}

func (c DataCreate) CreateData(ctx context.Context) (*xdrbuild.Transaction, error) {
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

	return tx, nil
}

// from https://gitlab.com/tokend/api/

// func checkTxResult(result horizon.SubmitResult) error {
// 	switch result.TXCode {
// 	case "tx_success":
// 		return nil
// 	case "tx_failed":
// 		return validation.Errors{"create account tx": errors.New(fmt.Sprint("transaction failed with op codes: ", result.OpCodes))}
// 	default:
// 		return validation.Errors{"create account tx": errors.New(fmt.Sprint("transaction failed with tx code: ", result.TXCode))}
// 	}
// }

// func convertTXSubmitError(s horizon.SubmitResult) (err error) {
// 	badRequests := map[string]struct{}{
// 		"op_invalid_destination": {},
// 		"op_already_exists":      {},
// 	}

// 	if len(s.OpCodes) > 0 {
// 		// only create account op codes are handled
// 		opCode := s.OpCodes[0]
// 		if _, ok := badRequests[opCode]; ok {
// 			return validation.Errors{
// 				"/data/attributes/account_id": errors.New(
// 					fmt.Sprintf("'%s' op code received on create account op", opCode),
// 				),
// 			}
// 		}
// 	}

// 	return checkTxResult(s)
// }
