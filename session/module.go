package session

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newFX),
)

type Params struct {
	fx.In

	SecretKey string `name:"secretKey"`
}

func newFX(params Params) *Store {
	return NewStore(StoreOptions{
		Key: params.SecretKey,
	})
}
