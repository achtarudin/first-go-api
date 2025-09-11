package infra

import (
	"fmt"

	"github.com/spf13/viper"
)

type appConfig struct {
	viper *viper.Viper
}

func NewAppConfig() *appConfig {
	return &appConfig{
		viper: viper.New(),
	}
}

func (a *appConfig) LoadEnvConfig(filename *string) error {

	envFile := ".env"

	if filename != nil && *filename != "" {
		envFile = *filename
	}

	a.viper.SetConfigFile(envFile)

	if err := a.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

func (a *appConfig) LoadTranslationConfig(filename *string) error {

	langFile := "lang/id.json"

	if filename != nil && *filename != "" {
		langFile = *filename
	}

	a.viper.SetConfigFile(langFile)
	a.viper.SetConfigType("json")

	if err := a.viper.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to merge lang config file: %w", err)
	}

	return nil
}

func (a *appConfig) GetViper() *viper.Viper {
	return a.viper
}

func LoadConfig(filename *string) (config *viper.Viper, err error) {

	v := viper.New()

	envFile := ".env"

	if filename != nil {
		envFile = *filename
	}

	v.SetConfigFile(envFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return v, nil
}

func LoadConfigLang() (config *viper.Viper, err error) {

	v := viper.New()

	v.SetConfigFile("lang/id.json")
	v.SetConfigType("json")
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config file: %w", err)
	}

	return v, nil
}
