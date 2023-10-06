package models

import (
	"database/sql"
	"errors"
	"fmt"

	data "BlobApi/internal/data"
)

type BlobModel struct {
	DB *sql.DB
}

func (m *BlobModel) Insert(userID int, data string) (int, error) {
	query := `
	INSERT INTO my_table (user_id, data) 
	VALUES ($1, $2)
	RETURNING index;
	`

	var id int
	res, err := m.DB.Exec(query, userID, data)
	if err != nil {
		return 0, err
	}

	if cnt, err := res.RowsAffected(); err != nil || cnt != 1 {
		// Handling an error or situation when the number of affected rows is not equal to the expected number of rows
		return 0, fmt.Errorf("unexpected number of affected rows: %d", cnt)
	}

	// Additional query for ID retrieval
	err = m.DB.QueryRow("SELECT lastval();").Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *BlobModel) Get(id int) (*data.Blob, error) {
	query := `
	SELECT index, user_id, data
	FROM my_table
	WHERE index = $1;
	`

	b := &data.Blob{}
	err := m.DB.QueryRow(query, id).Scan(&b.ID, &b.UserID, &b.Data)
	if err == sql.ErrNoRows {
		return nil, errors.New("blob not found")
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

func (m *BlobModel) GetBlobList() ([]*data.Blob, error) {
	query := `
	SELECT index, user_id, data
	FROM my_table;
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blobs []*data.Blob
	for rows.Next() {
		b := &data.Blob{}
		err := rows.Scan(&b.ID, &b.UserID, &b.Data)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blobs, nil
}

func (m *BlobModel) Delete(id int) error {
	query := `
	DELETE FROM my_table
	WHERE index = $1;
	`

	result, err := m.DB.Exec(query, id)
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

