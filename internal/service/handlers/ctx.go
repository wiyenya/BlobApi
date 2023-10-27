package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	//"gitlab.com/tokend/keypair"
)

type ctxKey int

const (
	logCtxKey       ctxKey = iota
	txBuilderCtxKey        = 4
	coreInfoCtxKey         = 8
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

// func DataCreate(r *http.Request) DataCreate {

// 	txbuilderbuilder := r.Context().Value(txBuilderCtxKey).(data.Infobuilder)
// 	info := r.Context().Value(coreInfoCtxKey).(data.Info)
// 	master := keypair.MustParseAddress(CoreInfo(r).MasterAccountID)
// 	tx := txbuilderbuilder(info, master)

// 	return accountcreator.New(
// 		tx,
// 		Horizon(r),
// 	)
// }
