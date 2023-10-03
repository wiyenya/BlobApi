package service

import (
	postgres "BlobApi/internal/data/postgres"
	"BlobApi/internal/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"

	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func (s *service) router() chi.Router {

	// Initializing the database
	connStr := "postgres://postgres:kate123@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() // Ping the database to make sure the connection is established
	if err != nil {
		log.Fatal(err)
	}

	//	Create BlobModel instance
	BlobModel := &postgres.BlobModel{DB: db}

	//	Create BlobHandler instance
	handler := handlers.NewBlobHandler(BlobModel)

	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)

	r.Route("/integrations/BlobApi", func(r chi.Router) {
		r.Post("/", handler.CreateBlob)
		r.Get("/", handler.GetBlobList)           // Получение списка блобов
		r.Get("/{blobID}", handler.GetBlobID)     // Получение блоба по ID
		r.Delete("/{blobID}", handler.DeleteBlob) // Удаление блоба по ID
	})

	return r
}
