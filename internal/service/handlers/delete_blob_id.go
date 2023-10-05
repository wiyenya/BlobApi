package handlers

import (
	"net/http"

	requests "BlobApi/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func (h *BlobHandler) DeleteBlob(w http.ResponseWriter, r *http.Request) {

	id, err := requests.DecodeDeleteBlobRequest(r)
	if err != nil || id < 1 {

		Log(r).WithError(err).Error("Invalid ID")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	error := h.Model.Delete(id)
	if error != nil {

		Log(r).WithError(err).Error("error deleting blob:")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
