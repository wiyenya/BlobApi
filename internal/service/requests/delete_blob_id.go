package requests

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func DecodeDeleteBlobRequest(r *http.Request) (int, error) {

	idStr := chi.URLParam(r, "blob_id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, err
	}

	return id, nil
}
