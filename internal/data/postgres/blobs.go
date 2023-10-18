package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	data "BlobApi/internal/data"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type BlobModel struct {
	DB *pgdb.DB
}

func (m *BlobModel) Insert(userID int, data map[string]interface{}) (int, error) {

	//data to JSON (map to bytes)
	jsonData, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return 0, errMarshal
	}

	// Using Squirrel to build an SQL query
	insertBuilder := sq.Insert("my_table").
		Columns("user_id", "data").
		Values(userID, jsonData).
		Suffix("RETURNING index").
		PlaceholderFormat(sq.Dollar)

	var id int
	errExec := m.DB.Exec(insertBuilder)
	if errExec != nil {
		return 0, errExec
	}

	// Additional query for ID retrieval
	//errQueryRow := m.DB.QueryRow("SELECT lastval();").Scan(&id)
	errQueryRow := m.DB.Get(&id, insertBuilder)
	if errQueryRow != nil {

		fmt.Print(errQueryRow)
		return 0, errQueryRow
	}

	return id, nil
}

func (m *BlobModel) Get(id int) (*data.Blob2, error) {

	// Using Squirrel to build an SQL query
	getBuilder := sq.Select("index", "user_id", "data").
		From("my_table").
		Where(sq.Eq{"index": id}).
		PlaceholderFormat(sq.Dollar)

	var blob data.Blob

	errQueryRow := m.DB.Get(&blob, getBuilder)
	if errQueryRow == sql.ErrNoRows {
		return nil, errors.New("blob not found")
	} else if errQueryRow != nil {
		return nil, errQueryRow
	}

	//JSON to Data (bytes to map)
	m1 := make(map[string]interface{})
	errUnmarshal := json.Unmarshal(blob.Data, &m1)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	Blob2 := &data.Blob2{}
	Blob2.Index = blob.Index
	Blob2.User_id = blob.User_id
	Blob2.Data = m1

	return Blob2, nil
}

// func (m *BlobModel) GetBlobList() ([]*data.Blob2, error) {
// 	// Using Squirrel to build an SQL query
// 	getBlobListBuilder := sq.Select("index", "user_id", "data").
// 		From("my_table").
// 		PlaceholderFormat(sq.Dollar)

// 	query, _, errGetBlobListBuilder := getBlobListBuilder.ToSql()
// 	if errGetBlobListBuilder != nil {
// 		return nil, errGetBlobListBuilder
// 	}

// 	rows, err := m.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var blobs []*data.Blob2
// 	for rows.Next() {
// 		b := &data.Blob{}
// 		err := rows.Scan(&b.ID, &b.UserID, &b.Data)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// JSON to Data (bytes to map)

// 		m1 := make(map[string]interface{})
// 		err1 := json.Unmarshal(b.Data, &m1)
// 		if err1 != nil {
// 			return nil, err
// 		}

// 		blob := &data.Blob2{}
// 		blob.ID = b.ID
// 		blob.UserID = b.UserID
// 		blob.Data = m1

// 		blobs = append(blobs, blob)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return blobs, nil
// }

// func (m *BlobModel) Delete(id int) error {

// 	// Using Squirrel to build an SQL query
// 	deleteBuilder := sq.Delete("my_table").
// 		Where(sq.Eq{"index": id}).
// 		PlaceholderFormat(sq.Dollar)

// 	query, _, errDeleteBuilder := deleteBuilder.ToSql()
// 	if errDeleteBuilder != nil {
// 		return errDeleteBuilder
// 	}

// 	result, err := m.DB.Exec(query, id)
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
