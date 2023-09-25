package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
	// "gitlab.com/tokend/go/signcontrol"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/signcontrol"
	// "gitlab.com/tokend/go/signcontrol"
	//"gitlab.com/tokend/go/signcontrol"
	//"gitlab.com/tokend/go/signcontrol"
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

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// fixme: skip_signature_check doesn't affect this function
	role, err := getSignerRole(r, request.Relationships.Owner.Data.ID)
	if err != nil {
		switch role {
		case badRequestFictiveRole:
			Log(r).WithError(err).Warn("bad request")
			ape.RenderErr(w, &jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprint(http.StatusBadRequest),
				Detail: "Request signature was invalid in some way",
			})
		case unauthorizedFictiveRole:
			Log(r).
				WithError(err).
				Warn("not admin's or owner's signer: only master admin may create blob with different owner")
			ape.RenderErr(w, problems.NotAllowed())
		case forbiddenFictiveRole:
			Log(r).WithError(err).Warn("attempt to create blob with non-registered signer")
			ape.RenderErr(w, problems.Forbidden())
		default:
			Log(r).WithError(err).Error("failed to get signer role")
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	blob, err := requests.Blob(request, role)
	if err != nil {
		Log(r).WithError(err).Warn("invalid blob type")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"/data/type": errors.New("invalid blob type")})...)
		return
	}

	err = BlobQ(r).Transaction(func(blobs data.Blobs) error {
		if err := blobs.Create(blob); err != nil {
			return errors.Wrap(err, "failed to create blob")
		}
		return nil
	})
	if err != nil {
		// silencing error to make request idempotent
		if errors.Cause(err) != postgres.ErrBlobsConflict {
			Log(r).WithError(err).Error("failed to save blob")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}
	// FIXME 201 only when actually created
	w.WriteHeader(201)
	ape.Render(w, &response)
}

func getSignerRole(r *http.Request, ownerAddress string) (uint64, error) {
	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		return badRequestFictiveRole, errors.Wrap(err, "bad signature")
	}

	findSigner := func(_signer, address string) (uint64, error) {
		accountSigners, err := AccountQ(r).Signers(address)
		if err != nil {
			return internalErrorFictiveRole, errors.Wrap(err, "failed to get owner signers")
		}
		if accountSigners == nil {
			return forbiddenFictiveRole, errors.New("account '" + address + "' does not have any signers")
		}
		// here AccountID is a public_key equivalent (not account_id) in DB, so it is compared correctly
		for _, s := range accountSigners {
			// note: add s.Weight > 0, like in doorman.signerOf func?
			if s.AccountID == _signer {
				return s.Role, nil
			}
		}
		return unauthorizedFictiveRole, errors.New("signer not found in account signers")
	}

	result, err := findSigner(signer, CoreInfo(r).GetMasterAccountID())
	if err == nil {
		return result, nil
	}
	return findSigner(signer, ownerAddress)
}
