package handlers

import (
	"net/http"
	// "gitlab.com/tokend/go/signcontrol"
	// "gitlab.com/tokend/api/internal/api/handlers/requests"
	// "gitlab.com/tokend/api/internal/api/resources"
	// "gitlab.com/tokend/api/internal/data"
	// "gitlab.com/tokend/api/internal/data/postgres"
)

const (
	badRequestFictiveRole    = 400
	unauthorizedFictiveRole  = 401
	forbiddenFictiveRole     = 403
	internalErrorFictiveRole = 500
)

func DeleteBlobBlobId(w http.ResponseWriter, r *http.Request) {

}
