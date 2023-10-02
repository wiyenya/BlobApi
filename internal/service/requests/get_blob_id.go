package requests

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func DecodeGetBlobRequest(r *http.Request) (int, error) {

	idStr := chi.URLParam(r, "blobID")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, err
	}

	return id, nil
}
