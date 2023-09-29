package requests

import (
	"encoding/json"
	"net/http"

	resourses "BlobApi/resources"
)

func DecodeCreateBlobRequest(r *http.Request) (*resourses.BlobRequest, error) {

	var req resourses.BlobRequestResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &req.Data, nil
}
