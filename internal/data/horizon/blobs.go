package postgres

import (
	"bytes"
	"net/http"

	"github.com/jmoiron/sqlx/types"
	//"gitlab.com/distributed_lab/kit/pgdb"
	dataPkg "BlobApi/internal/data"

	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

// type BlobModel struct {
// 	DB *pgdb.DB
// }

// func (m *BlobModel) Insert(userID int32, data types.JSONText) (int, error) {

// 	// Using Squirrel to build an SQL query
// 	insertBuilder := sq.Insert("my_table").
// 		Columns("user_id", "data").
// 		Values(userID, data).
// 		Suffix("RETURNING index").
// 		PlaceholderFormat(sq.Dollar)

// 	//insertBuilder - adds to the database, id - get the id

// 	var id int
// 	errGet := m.DB.Get(&id, insertBuilder)
// 	if errGet != nil {
// 		return 0, errGet
// 	}

// 	return id, nil
// }

func Insert(userID int32, data types.JSONText) (*xdr.Operation, error) {

	txBuilder := xdrbuild.NewBuilder("test", int64(1234))

	blob := dataPkg.Blob{
		Index:  1,
		UserId: &userID,
		Data:   data,
	}

	createData := xdrbuild.CreateData{
		Type:  uint64(12345),
		Value: blob,
	}

	operation, err := createData.XDR()
	if err != nil {
		return nil, err
	}

	// Добавляем операцию к транзакции
	tx := txBuilder.Transaction().Op(createData)

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

	endpoint := "http://localhost:8000/_/api/"
	resp, err := http.Post(endpoint, "application/base64", bytes.NewBufferString(encodedSignedTransaction))
	// if err != nil {
	// 	return "", err
	// }
	defer resp.Body.Close()

	return operation, nil
}

// func (m *BlobModel) Get(id int) (*data.Blob, error) {

// 	// Using Squirrel to build an SQL query
// 	getBuilder := sq.Select("index", "user_id", "data").
// 		From("my_table").
// 		Where(sq.Eq{"index": id}).
// 		PlaceholderFormat(sq.Dollar)

// 	var blob data.Blob

// 	errQueryRow := m.DB.Get(&blob, getBuilder)
// 	if errQueryRow == sql.ErrNoRows {
// 		return nil, errors.New("blob not found")
// 	} else if errQueryRow != nil {
// 		return nil, errQueryRow
// 	}

// 	return &blob, nil
// }

// func (m *BlobModel) GetBlobList() ([]*data.Blob, error) {
// 	// Using Squirrel to build an SQL query
// 	getBlobListBuilder := sq.Select("index", "user_id", "data").
// 		From("my_table").
// 		PlaceholderFormat(sq.Dollar)

// 	var blobs []*data.Blob
// 	err := m.DB.Select(&blobs, getBlobListBuilder)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var blobs2 []*data.Blob
// 	for _, blob := range blobs {
// 		blobs2 = append(blobs2, blob)
// 	}

// 	return blobs2, nil
// }

// func (m *BlobModel) Delete(id int) error {

// 	// Using Squirrel to build an SQL query
// 	deleteBuilder := sq.Delete("my_table").
// 		Where(sq.Eq{"index": id}).
// 		PlaceholderFormat(sq.Dollar)

// 	result, err := m.DB.ExecWithResult(deleteBuilder)
// 	if err != nil {
// 		return err
// 	}

// 	// Check that at least one line has been deleted
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("no rows affected, blob might not exist")
// 	}

// 	return nil
// }
