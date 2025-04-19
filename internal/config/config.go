package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	OpenApiToken     string `env:"TI_OPENAI_API_KEY" env-default:""`
	DiscordToken     string `env:"TI_DISCORD_BOT_TOKEN" env-required:"true"`
	Proxy            string `env:"TI_PROXY" env-default:""`
	ChannelDashboard string `env:"TI_DASHBOARD_CHANNEL" env-required:""`
	SizeContext      int    `env:"TI_SIZE_CONTEXT" env-default:"100"`
}

func MustLoad() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
