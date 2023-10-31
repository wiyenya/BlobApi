package handlers

import (
	horizon "BlobApi/internal/data/horizon"
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	//"gitlab.com/tokend/keypair"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	txBuilderCtxKey
	coreInfoCtxKey
	horizonCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxHorizonConnector(entry *horizon.HorizonModel) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonCtxKey, entry)
	}
}

func HorizonConnector(r *http.Request) *horizon.HorizonModel {
	return r.Context().Value(horizonCtxKey).(*horizon.HorizonModel)
}
