package data

import (
	"gitlab.com/tokend/regources"
)

//go:generate mockery -case underscore -name Info

type Info interface {
	Info() (*regources.Info, error)
}
