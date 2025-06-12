package impl

import "go.uber.org/fx"

// FXModule represents a FX module for app service.
var FXModule = fx.Provide(
	NewItemService,
)
