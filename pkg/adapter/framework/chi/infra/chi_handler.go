package infra

import (
	"context"
	"examples/pkg/adapter/infra"
)

type ChiHandler func(ctx context.Context, chiCtx ChiContext) *infra.HandleError

type ChiContext interface {
	URLParam(name string) string
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, res *infra.HTTPError) error
}
