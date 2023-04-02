package repository

import "context"

type TxExecFunc func(ctx context.Context) (interface{}, error)

type Transaction interface {
	Do(ctx context.Context, fn TxExecFunc) (interface{}, error)
}
