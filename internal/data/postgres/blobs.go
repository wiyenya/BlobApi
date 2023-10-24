package postgres

import (
	"database/sql"
	"errors"

	data "BlobApi/internal/data"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type BlobModel struct {
	DB *pgdb.DB
}

func (m *BlobModel) Insert(userID int32, data []byte) (int, error) {

	// Using Squirrel to build an SQL query
	insertBuilder := sq.Insert("my_table").
		Columns("user_id", "data").
		Values(userID, data).
		Suffix("RETURNING index").
		PlaceholderFormat(sq.Dollar)

	//insertBuilder - adds to the database, id - get the id

	var id int
	errGet := m.DB.Get(&id, insertBuilder)
	if errGet != nil {
		return 0, errGet
	}

	return id, nil
}

func (m *BlobModel) Get(id int) (*data.Blob, error) {

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

	return &blob, nil
}

func (m *BlobModel) GetBlobList() ([]*data.Blob, error) {
	// Using Squirrel to build an SQL query
	getBlobListBuilder := sq.Select("index", "user_id", "data").
		From("my_table").
		PlaceholderFormat(sq.Dollar)

	var blobs []*data.Blob
	err := m.DB.Select(&blobs, getBlobListBuilder)
	if err != nil {
		return nil, err
	}

	var blobs2 []*data.Blob
	for _, blob := range blobs {
		blobs2 = append(blobs2, blob)
	}

	return blobs2, nil
}

func (m *BlobModel) Delete(id int) error {

	// Using Squirrel to build an SQL query
	deleteBuilder := sq.Delete("my_table").
		Where(sq.Eq{"index": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := m.DB.ExecWithResult(deleteBuilder)
	if err != nil {
		return err
	}

	// Check that at least one line has been deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected, blob might not exist")
	}

	return nil
}
