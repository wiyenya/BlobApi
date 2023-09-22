package service

import (
	"net"
	"net/http"

	"BlobApi/internal/config"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/api/internal/api/handlers"
    "gitlab.com/tokend/api/internal/api/middlewares"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}

func BlobsRouter(entry *logan.Entry) chi.Router {
    r := chi.NewRouter()

    r.Use(
        ape.RecoverMiddleware(entry),
        middlewares.Logger(entry, 300*time.Millisecond),
    )

	// blobs
    r.Route("/integrations/BlobApi", func(r chi.Router) {
        r.Post("/", handlers.CreateBlob)
        r.Get("/", handlers.GetBlobList)         // Получение списка блобов
        r.Get("/{blobID}", handlers.GetBlob)    // Получение блоба по ID
        r.Delete("/{blobID}", handlers.DeleteBlob) // Удаление блоба по ID
    })


    return r
}