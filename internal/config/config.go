package config

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Alert struct {
		CheckInterval time.Duration `mapstructure:"check_interval"`
	} `mapstructure:"alert"`

	Recorder struct {
		Extension       string   `mapstructure:"extension"`
		OutputTemplate  string   `mapstructure:"output_template"`
		CommandTemplate []string `mapstructure:"command_template"`
	} `mapstructure:"recorder"`

	Following struct {
		Nico []string `mapstructure:"nico"`
	} `mapstructure:"following"`

	Auth struct {
		Nico struct {
			UserSession string `mapstructure:"user_session"`
		} `mapstructure:"nico"`
	} `mapstructure:"auth"`

	WebPush struct {
		NicoPush struct {
			UAID       string   `mapstructure:"uaid"`
			AuthSecret string   `mapstructure:"auth_secret"`
			PrivateKey string   `mapstructure:"private_key"`
			ChannelIDs []string `mapstructure:"channel_ids"`
		} `mapstructure:"nicopush"`
	} `mapstructure:"webpush"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	viper.SetConfigFile(path)

	// Default values
	viper.SetDefault("alert.check_interval", "10s")
	viper.SetDefault("recorder.extension", "ts")
	viper.SetDefault("recorder.output_template",
		`{{.StartedAt.Format "20060102"}}-{{.ProgramId}}-{{printf "%.20s" .AuthorName}}-{{printf "%.50s" .ProgramTitle}}`,
	)
	viper.SetDefault("recorder.command_template", []string{
		"ffmpeg",
		"-cookies", `{{formatCookies .Cookies "\n"}}`,
		"-i", "{{.URL}}",
		"-c", "copy",
		"{{.Output}}.{{.Extension}}",
	})

	load := func() {
		if err := viper.Unmarshal(&config); err != nil {
			log.Println(err)
		}
	}

	// Load config from path
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
	load()

	// Live Reload
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config: file changed", e.Name)
		load()
	})
	viper.WatchConfig()

	return &config, nil
}
