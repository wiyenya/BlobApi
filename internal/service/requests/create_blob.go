package requests

import (
	"encoding/json"
	"net/http"

	resourses "BlobApi/resources"

	"github.com/pkg/errors"
)

func DecodeCreateBlobRequest(r *http.Request) (resourses.BlobRequest, error) {

	var req resourses.BlobRequestResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req.Data, errors.Wrap(err, "failed to unmarshal")
	}

	return req.Data, nil
}
