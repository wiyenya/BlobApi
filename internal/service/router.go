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

	//инициализация базы данных
	connStr := "postgres://postgres:kate123@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() // Пинговать базу данных, чтобы убедиться, что соединение установлено
	if err != nil {
		log.Fatal(err)
	}

	//Создание Экземпляра BlobModel
	yourBlobModel := &postgres.BlobModel{DB: db}

	//Создание Экземпляра BlobHandler
	handler := handlers.NewBlobHandler(yourBlobModel)

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
		r.Get("/", handlers.GetBlobList)           // Получение списка блобов
		r.Get("/{blobID}", handlers.GetBlobID)     // Получение блоба по ID
		r.Delete("/{blobID}", handlers.DeleteBlob) // Удаление блоба по ID
	})

	return r
}
