package service

import (
	postgres "BlobApi/internal/data/postgres"
	"BlobApi/internal/service/handlers"
	middleware "BlobApi/internal/service/middleware"
	"errors"
	"os"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"

	"log"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/figure"
)

type Config struct {
	DBConfig struct {
		DatabaseURL string `fig:"database_url,required"`
	} `fig:"config"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	// Read the value of the environment variable
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		return nil, errors.New("environment variable KV_VIPER_FILE is not set")
	}

	v.SetConfigFile(configPath)

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

func (s *service) router(entry *logan.Entry) chi.Router {

	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbOpts := pgdb.Opts{
		URL:                cfg.DBConfig.DatabaseURL,
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
