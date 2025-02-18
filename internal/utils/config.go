package utils

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Alert struct {
		CheckIntervalSec time.Duration `mapstructure:"check_interval_sec"`
	} `mapstructure:"alert"`
	Recorder struct {
		OutputTemplate  string   `mapstructure:"output_template"`
		CommandTemplate []string `mapstructure:"command_template"`
	} `mapstructure:"recorder"`
	Following struct {
		Nico []string `mapstructure:"nico"`
	} `mapstructure:"following"`
}

func LoadConfig(path string) (*Config, error) {
	// Default values
	viper.SetDefault("alert.check_interval_sec", 10)
	viper.SetDefault("recorder.output_template",
		"{yyyymmdd}-{id}-{providerId}-{title}",
	)
	viper.SetDefault("recorder.command_template", []string{
		"ffmpeg",
		"-cookies", "{cookies}",
		"-i", "{url}",
		"-c", "copy",
		"{output}",
	})

	// Load config from path
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Panic(err)
	}

	return config, nil
}
