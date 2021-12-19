package templates

import (
	"emailaddress.horse/thousand/session"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fxNewRenderer),
)

type RendererParams struct {
	fx.In

	Store *session.Store
}

func fxNewRenderer(params RendererParams) *Renderer {
	return NewRenderer(RendererOptions{
		Store: params.Store,
	})
}
