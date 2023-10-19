package service

import (
	postgres "BlobApi/internal/data/postgres"
	"BlobApi/internal/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/pgdb"

	"log"

	_ "github.com/lib/pq"


	"gitlab.com/distributed_lab/figure"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `fig:"database_url,required"`
}

func NewConfig(filePath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filePath)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	err := figure.Out(cfg).From(v.AllSettings()).Please()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *service) router() chi.Router {

	cfg, err := NewConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbOpts := pgdb.Opts{
		URL:                cfg.DatabaseURL,
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
	}

	// Open a connection to the database
	db, err := pgdb.Open(dbOpts)
	if err != nil {
		log.Fatal(err)
	}

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
	)

	r.Route("/integrations/blobs", func(r chi.Router) {
		r.Post("/", handler.CreateBlob)
		r.Get("/", handler.GetBlobList)
		r.Get("/{blob_id}", handler.GetBlobID)
		r.Delete("/{blob_id}", handler.DeleteBlob)
	})

	return r
}
