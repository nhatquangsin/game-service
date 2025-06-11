package repoimpl

import (
	"go.uber.org/fx"
)

// FXModule represents a FX module for options.
var FXModule = fx.Provide(
	NewItemRepo,
)
