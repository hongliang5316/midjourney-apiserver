package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort int32      `yaml:"listen_port"`
	Midjourney Midjourney `yaml:"midjourney"`
	Redis      Redis      `yaml:"redis"`
}

type Midjourney struct {
	UserToken string `yaml:"user_token"`
	GuildID   string `yaml:"guild_id"`
	ChannelID string `yaml:"channel_id"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func Load() *Config {
	cfg := new(Config)

	data, err := ioutil.ReadFile("./conf/conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.ListenPort == 0 {
		cfg.ListenPort = 8080
	}

	return cfg
}
