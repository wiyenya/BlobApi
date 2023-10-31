package service

import (
	"BlobApi/internal/config"
	//postgres "BlobApi/internal/data/postgres"
	horizon "BlobApi/internal/data/horizon"
	"BlobApi/internal/service/handlers"
	middleware "BlobApi/internal/service/middleware"

	"time"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

func (s *service) router(entry *logan.Entry, cfg config.Config) chi.Router {

	connector := horizon.NewHorizonModel(entry, "http://localhost:8000/_/api/",
		"SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")

	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxHorizonConnector(connector),
		),
		middleware.Logger(entry, 300*time.Millisecond),
	)

	r.Route("/integrations/blobs", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/", handlers.GetBlobList)
		r.Get("/{blob_id}", handlers.GetBlobID)
		r.Delete("/{blob_id}", handlers.DeleteBlob)
	})

	return r
}
