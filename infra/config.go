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

	a.viper.AutomaticEnv()

	// envFile := ".env"

	// if filename != nil && *filename != "" {
	// 	envFile = *filename
	// }

	// a.viper.SetConfigFile(envFile)

	// if err := a.viper.ReadInConfig(); err != nil {
	// 	return fmt.Errorf("failed to read config file: %w", err)
	// }

	return nil
}

func (a *appConfig) LoadTranslationConfig(filename *string) error {

	langFile := "lang/id.json"

	if filename != nil && *filename != "" {
		langFile = *filename
	}

	_ = langFile
	// a.viper.SetConfigFile(langFile)
	// a.viper.SetConfigType("json")

	a.viper.SetConfigName("id")         // path to look for the config file in
	a.viper.SetConfigType("json")       // path to look for the config file in
	a.viper.AddConfigPath("/app/lang/") // path to look for the config file in
	a.viper.AddConfigPath(".")          // optionally look for config in the working directory

	if err := a.viper.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to merge lang config file: %w", err)
	}

	return nil
}

func (a *appConfig) GetViper() *viper.Viper {
	return a.viper
}
