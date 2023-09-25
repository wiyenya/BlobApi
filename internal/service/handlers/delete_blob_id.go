package handlers

import (
	"net/http"
)

const (
	badRequestFictiveRole    = 400
	unauthorizedFictiveRole  = 401
	forbiddenFictiveRole     = 403
	internalErrorFictiveRole = 500
)

func DeleteBlobBlob(w http.ResponseWriter, r *http.Request) {

}
