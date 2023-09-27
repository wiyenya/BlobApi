package service

import (
	"BlobApi/internal/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)

	r.Route("/integrations/BlobApi", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/", handlers.GetBlobID)             // Получение списка блобов
		r.Get("/{blobID}", handlers.GetBlob)       // Получение блоба по ID
		r.Delete("/{blobID}", handlers.DeleteBlob) // Удаление блоба по ID
	})

	return r
}
