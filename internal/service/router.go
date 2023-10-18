package service

import (
	postgres "BlobApi/internal/data/postgres"
	"BlobApi/internal/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/pgdb"

	"log"

	_ "github.com/lib/pq"
)

func (s *service) router() chi.Router {

	// Создаем настройки подключения к базе данных
	dbOpts := pgdb.Opts{
		URL:                "postgres://postgres:kate123@localhost:5432/mydatabase?sslmode=disable",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
	}

	// Открываем соединение с базой данных
	db, err := pgdb.Open(dbOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем экземпляр BlobModel, передавая ему db
	BlobModel := &postgres.BlobModel{DB: db}

	// Создаем экземпляр BlobHandler, передавая ему BlobModel
	handler := handlers.NewBlobHandler(BlobModel)

	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)

	r.Route("/integrations/blobs", func(r chi.Router) {
		r.Post("/", handler.CreateBlob)
		//r.Get("/", handler.GetBlobList)
		r.Get("/{blob_id}", handler.GetBlobID)
		r.Delete("/{blob_id}", handler.DeleteBlob)
	})

	return r
}
