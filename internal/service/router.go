package service

import (
	"BlobApi/internal/config"
	"context"
	"sync"

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

	// Create an instance of BlobModel by passing db to it
	connector := horizon.NewHorizonModel(entry, "http://localhost:8000/_/api/",
		"SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")

	// Create an instance of BlobHandler by passing it a BlobModel
	handler := handlers.NewBlobHandler(BlobModel)

	sync.WaitGroup

	r := chi.NewRouter()

	ctx := context.Background()
	ctx, _ := context.WithDeadline(ctx, time.Now().Add(30*time.Millisecond))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Do something
			}
		}
	}()

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
		r.Post("/", handler.CreateBlob)
		r.Get("/", handler.GetBlobList)
		r.Get("/{blob_id}", handler.GetBlobID)
		r.Delete("/{blob_id}", handler.DeleteBlob)
	})

	return r
}
