package requests

import (
	"encoding/json"
	"net/http"
)

type CreateBlobRequest struct {
	UserID int    `json:"user_id"`
	Data   string `json:"data"`
}

func DecodeCreateBlobRequest(r *http.Request) (*CreateBlobRequest, error) {
	var req CreateBlobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
