package session

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newFX),
)

func newFX() *Store {
	return NewStore(StoreOptions{
		Key: "secret",
	})
}
