package viperutil

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stoewer/go-strcase"
	"go.uber.org/fx"
)

const (
	envPrefix = "kit"

	envKeyEnv = "env"
	// envKeyEnvConfigOnly = "env_config_only"
	envKeyEnvConfigOnly = "ENV_CONFIG_ONLY"
	envKeyConfigPath    = "kit_config_path"

	envDefaultValEnv = "dev"
)

// FXModule represents a FX module for Viper.
var FXModule = fx.Options(
	fx.Provide(
		NewViper,
	),
)

// NewViper creates a new instance of Viper with default path "./config".
//
// The configuration file is named by "ENV" variables and in "yaml" format.
//
//   - ENV=dev 			=> configFile=config/config.dev.yaml
//   - ENV=development 	=> configFile=config/config.development.yaml
//
// Example:
//
//	viper, err := viperutil.NewViper()
//	if err != nil {
//		log.Println(err)
//	}
func NewViper() (*viper.Viper, error) {
	configPath := os.Getenv(envKeyConfigPath)
	if configPath == "" {
		configPath = "./config"
	}
	return NewViperFrom(configPath)
}

// NewViperFrom creates a new instance of Viper with input path.
//
// The configuration file is named by "ENV" variables and in "yaml" format.
//
//   - ENV=dev 			=> configPath=config/config.dev.yaml
//   - ENV=development 	=> configPath=config/config.development.yaml
//
// Example:
//
//	viper, err := viperutil.NewViperFrom("./config/path")
//	if err != nil {
//		log.Println(err)
//	}
func NewViperFrom(path string) (*viper.Viper, error) {
	v := viper.New()

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	v.SetDefault(envKeyEnv, envDefaultValEnv)
	env := v.GetString(envKeyEnv)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//envOnly := v.GetBool(envKeyEnvConfigOnly)
	isEnvConfig := os.Getenv(envKeyEnvConfigOnly)
	if isEnvConfig == "" {
		isEnvConfig = "false"
	}
	envOnly, err := strconv.ParseBool(isEnvConfig)
	if err != nil {
		return nil, err
	}

	if !envOnly {
		v.SetConfigName("config." + env)
		v.SetConfigType("yaml")
		v.AddConfigPath(path)

		return v, v.ReadInConfig()
	}

	return v, nil
}

// StringToSliceHookFunc returns a DecodeHookFunc that converts
// string to []string by splitting on the given sep.
func StringToSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Kind,
		t reflect.Kind,
		data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Slice {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		return strings.Split(strings.Trim(raw, "\n"), sep), nil
	}
}

// StringToByteSliceHookFunc returns a DecodeHookFunc that converts
// base64 string to []byte.
func StringToByteSliceHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		sDec, err := base64.StdEncoding.DecodeString(data.(string))
		if err != nil {
			return nil, err
		}

		return sDec, nil
	}
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
// Unmarshal use
//
//	 mapstructure.ComposeDecodeHookFunc(
//			mapstructure.StringToTimeDurationHookFunc(),
//			viperutil.StringToSliceHookFunc(","),
//		)
//
// as default decode hook funcs.
func Unmarshal(v *viper.Viper, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if err := bindStructPathsToEnvs(v, rawVal); err != nil {
		return err
	}

	return v.Unmarshal(
		rawVal,
		append(
			opts,
			viper.DecodeHook(
				mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					StringToSliceHookFunc(" "),
					StringToByteSliceHookFunc(),
				),
			),
		)...,
	)
}

// bindStructPathsToEnvs binds all structure paths extract from tag
// "mapstructure" to environment variables.
func bindStructPathsToEnvs(v *viper.Viper, rawVal interface{}) error {
	paths, err := extractMapstructurePaths("", reflect.TypeOf(rawVal))
	if err != nil {
		return err
	}

	envNames := getAllEnvNamesWithPrefix(envPrefix)
	for _, path := range paths {
		if strings.IndexByte(path, byte('*')) == -1 {
			if err := v.BindEnv(path); err != nil {
				return err
			}
		}

		pathFmt := strings.Replace(path, "*", "%s", -1)
		envNameRegex := regexp.MustCompile(strings.Replace(strings.Replace(path, ".", "_", -1), "*", "(.*)", -1))
		for _, envName := range envNames {
			if envNameRegex.MatchString(envName) {
				matches := envNameRegex.FindStringSubmatch(envName)
				if len(matches) > 1 {
					var matchesIface []interface{}
					for _, m := range matches[1:] {
						matchesIface = append(matchesIface, m)
					}
					if err := v.BindEnv(fmt.Sprintf(pathFmt, matchesIface...)); err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

// extractMapstructurePaths extracts struct path from the input type.
func extractMapstructurePaths(parent string, typ reflect.Type) ([]string, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var res []string
	switch typ.Kind() {
	case reflect.Map:
		if typ.Key().Kind() == reflect.String {
			mapElemTyp := typ.Elem()
			for mapElemTyp.Kind() == reflect.Ptr {
				mapElemTyp = mapElemTyp.Elem()
			}

			switch mapElemTyp.Kind() {
			case reflect.Struct:
				childMapPaths, err := extractMapstructurePaths(parent+".*", mapElemTyp)
				if err != nil {
					return nil, err
				}
				res = append(res, childMapPaths...)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
				res = append(res, parent+".*")
			default:

			}
		}
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldTyp := field.Type

			path, err := getMapstructurePath(parent, field)
			if err != nil {
				return nil, err
			}

			if field.Anonymous {
				paths, err := extractMapstructurePaths(parent, fieldTyp)
				if err != nil {
					return nil, err
				}
				res = append(res, paths...)
			} else if fieldTyp.Kind() == reflect.Ptr || fieldTyp.Kind() == reflect.Struct || fieldTyp.Kind() == reflect.Map {
				paths, err := extractMapstructurePaths(path, fieldTyp)
				if err != nil {
					return nil, err
				}
				res = append(res, paths...)
			} else {
				res = append(res, path)
			}
		}
	default:
		res = append(res, parent)
	}

	return res, nil
}

// getMapstructurePath returns the structure path of the input field.
func getMapstructurePath(parent string, field reflect.StructField) (string, error) {
	path := strcase.LowerCamelCase(field.Name)
	if parent != "" {
		path = parent + "." + path
	}

	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return "", err
	}

	tag, err := tags.Get("mapstructure")
	if err != nil {
		return path, nil
	}

	path = strings.ToLower(tag.Name)
	if parent != "" {
		path = parent + "." + path
	}

	return path, nil
}

// getAllEnv returns all name of environment variables in lower case.
func getAllEnvNamesWithPrefix(prefix string) []string {
	var envNames []string

	for _, n := range os.Environ() {
		envElems := strings.Split(n, "=")
		if len(envElems) > 0 {
			envName := strings.ToLower(envElems[0])
			if strings.HasPrefix(envName, prefix) {
				envNames = append(envNames, strings.TrimPrefix(envName, prefix+"_"))
			}
		}
	}

	return envNames
}
