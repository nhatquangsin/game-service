package config

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"

	viperutil "github.com/nhatquangsin/game-service/infra/utils/viper"
)

// Config is a group of options for the service.
type Config struct {
	Environment string `mapstructure:"env"`

	Service struct {
		Name      string `json:"name" mapstructure:"name"`
		Component string `json:"component" mapstructure:"component"`
		Namespace string `json:"namespace" mapstructure:"namespace"`
		Version   string `json:"version" mapstructure:"version"`
	} `mapstructure:"service"`
}

// Load loads Config from Viper and returns them.
func Load(v *viper.Viper) (Config, error) {
	cfg := Config{}
	if err := viperutil.Unmarshal(v, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// FXModule represents a FX module for config.
var FXModule = fx.Options(
	fx.Provide(
		Load,
	),
)
