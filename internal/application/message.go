package application

import "github.com/bwmarrin/discordgo"

func (app *Application) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != app.Cfg.ChannelID {
		return
	}
}

func (app *Application) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.ChannelID != app.Cfg.ChannelID {
		return
	}
}
