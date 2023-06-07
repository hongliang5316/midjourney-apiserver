package application

type Config struct {
	UserToken  string `yaml:"user_token"`
	GuildID    string `yaml:"guild_id"`
	ChannelID  string `yaml:"channel_id"`
	ListenPort int32  `yaml:"listen_port"`
}
