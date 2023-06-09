package common

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hongliang5316/midjourney-apiserver/internal/config"
	"github.com/hongliang5316/midjourney-apiserver/internal/store"
	"github.com/hongliang5316/midjourney-go/midjourney"
)

type Base struct {
	*discordgo.Session
	Store    *store.Store
	MJClient *midjourney.Client
	Config   *config.Config
}
