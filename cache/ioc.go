package cache

import "go.uber.org/fx"

// FXModule represents a FX module for cache service.
var FXModule = fx.Provide(
	LoadAllItems,
)
