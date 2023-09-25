package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/doorman"
	// "gitlab.com/tokend/api/internal/api/resources"
	// "gitlab.com/tokend/api/internal/types"
	// "gitlab.com/tokend/go/doorman"
)

type (
	GetBlobRequest struct {
		Address types.Address `json:"-"`
		BlobID  string        `json:"-"`
	}
)

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, error) {
	request := GetBlobRequest{
		Address: types.Address(chi.URLParam(r, "address")),
		BlobID:  chi.URLParam(r, "blob"),
	}
	return request, request.Validate()
}

func (r GetBlobRequest) Validate() error {
	err := Errors{
		"blob": Validate(&r.BlobID, Required),
	}
	return err.Filter()
}

func GetBlob(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("tag", "blob_performance")
	log.Info("Request started")
	request, err := NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log.Info("Try to get blob from DB")
	blob, err := BlobQ(r).Get(request.BlobID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !types.IsPublicBlob(blob.Type) {
		log.Info("Checking permissions")
		constrains := []doorman.SignerConstraint{doorman.SignerOf(CoreInfo(r).GetMasterAccountID())}

		if blob.OwnerAddress != nil {
			ownerAddress := blob.OwnerAddress.String()

			log.Info("Get KYC recovery signer role")
			kycRole, err := SystemSettings(r).KYCRecoverySignerRole()
			if err != nil {
				Log(r).WithError(err).Error("failed to get kyc recovery signer role")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			// Default case:
			//  not setting any constrains here because doorman is initialized
			//  with restricted `kycRecoverySignerRole` and `licenseAdminSignerRole`
			//  so only signers with roles different from mentioned above can access the blob
			constrain := doorman.SignerOf(ownerAddress)

			if blob.CreatorSignerRole != nil && *blob.CreatorSignerRole == kycRole {
				constrain = doorman.ClearSignerOf(ownerAddress)
			}

			constrains = append(constrains, constrain)
		}

		log.Info("Check permissions with doorman")
		if err := Doorman(r, constrains...); err != nil {
			ape.RenderErr(w, problems.NotAllowed(err))
			return
		}
	}

	log.Info("Render response")
	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}

	ape.Render(w, &response)
}

func NewBlob(blob *types.Blob) resources.Blob {
	b := resources.Blob{
		Key: resources.Key{
			ID:   blob.ID,
			Type: resources.ResourceType(blob.Type.String()),
		},
		Attributes: resources.BlobAttributes{
			Value: blob.Value,
		},
	}
	return b
}
