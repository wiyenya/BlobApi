package service

import (
	"BlobApi/internal/config"
	postgres "BlobApi/internal/data/postgres"
	"BlobApi/internal/service/handlers"
	middleware "BlobApi/internal/service/middleware"

	"time"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

func (s *service) router(entry *logan.Entry, cfg config.Config) chi.Router {

	// Open a connection to the database
	db := cfg.DB()

	// Create an instance of BlobModel by passing db to it
	BlobModel := &postgres.BlobModel{DB: db}

	// Create an instance of BlobHandler by passing it a BlobModel
	handler := handlers.NewBlobHandler(BlobModel)

	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
		middleware.Logger(entry, 300*time.Millisecond),
	)

	r.Route("/integrations/blobs", func(r chi.Router) {
		r.Post("/", handler.CreateBlob)
		r.Get("/", handler.GetBlobList)
		r.Get("/{blob_id}", handler.GetBlobID)
		r.Delete("/{blob_id}", handler.DeleteBlob)
	})

	return r
}
