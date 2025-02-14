package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Meta struct {
		FetchInterval time.Duration `mapstructure:"fetch_interval"`
	} `mapstructure:"meta"`

	Paths struct {
		FFmpeg        string `mapstructure:"ffmpeg"`
		OutputBaseDir string `mapstructure:"output_base_dir"`
	} `mapstructure:"paths"`

	Following struct {
		Users map[string][]string `mapstructure:"users"`
	} `mapstructure:"following"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.SetDefault("meta.fetch_interval", 10)
	viper.SetDefault("paths.ffmpeg", "ffmpeg")
	viper.SetDefault("paths.output_base_dir", ".")

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to unmarshal struct: %v", err)
	}

	return config, nil
}
