package config

import (
	"encoding/json"
	"sync"

	"github.com/spf13/viper"
	"github.com/yimincai/health-checker/pkg/logger"
)

type Config struct {
	Prefix               string `mapstructure:"PREFIX"`
	DiscordToken         string `mapstructure:"DISCORD_TOKEN"`
	Endpoint             string `mapstructure:"ENDPOINT"`
	PunchApiPath         string `mapstructure:"PUNCH_API_PATH"`
	LoginApiPath         string `mapstructure:"LOGIN_API_PATH"`
	NotificationChannnel string `mapstructure:"NOTIFICATION_CHANNEL"`
}

var cfg *Config
var cfgOnce sync.Once

// New init env, this function will load .env file at first if exist it will load environment variable
// APP_ENV is required, if not set, it will panic
// .env file is for local development, environment variable is for production
func New() *Config {
	cfgOnce.Do(func() {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Warnf("Using environment variables")

			viper.AutomaticEnv()
			_ = viper.BindEnv("PREFIX")
			_ = viper.BindEnv("DISCORD_TOKEN")
			_ = viper.BindEnv("ENDPOINT")
			_ = viper.BindEnv("PUNCH_API_PATH")
			_ = viper.BindEnv("LOGIN_API_PATH")
			_ = viper.BindEnv("SENDING_CHANNEL")
			_ = viper.BindEnv("NOTIFICATION_CHANNEL")
		}

		err = viper.Unmarshal(&cfg)
		if err != nil {
			logger.Fatalf("Environment can't be loaded: %s", err)
		}

		if cfg.DiscordToken == "" {
			panic("Discord token is required")
		}

		logger.Infof("Config: \n%s", prettyPrint(cfg))
	})

	return cfg
}

func GetEnv() *Config {
	return cfg
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
