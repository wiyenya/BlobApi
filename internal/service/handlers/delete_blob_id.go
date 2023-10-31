package handlers

import (
	"net/http"

	requests "BlobApi/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {

	id, err := requests.DecodeDeleteBlobRequest(r)
	if err != nil || id < 1 {

		Log(r).WithError(err).Error("Invalid ID")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	connector := HorizonConnector(r)

	errorDelete := connector.Delete(id)
	if errorDelete != nil {
		Log(r).WithError(errorDelete).Error("error deleting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
