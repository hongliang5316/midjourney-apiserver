package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hongliang5316/midjourney-apiserver/pkg/store"
	wb "github.com/hongliang5316/midjourney-apiserver/pkg/webhook"
)

func (app *Application) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != app.Config.Midjourney.ChannelID {
		return
	}

	if m.Author.Username == "Midjourney Bot" {
		log.Printf("%s, messageCreate: %s", m.ID, m.Content)
		// log.Printf("%s, messageCreate: %s", m.ID, toJson(m))

		if m.Interaction != nil && m.Interaction.Name == "describe" {
			app.handleDescribeEvent(m)
			return
		}

		// handle wating to start message
		if strings.HasSuffix(m.Content, "(Waiting to start)") {
			app.handleWaitingToStartEvent(m)
			return
		}

		// handle embed error message
		if len(m.Attachments) == 0 && len(m.Embeds) > 0 {
			app.handleEmbedErrorEvent(m)
			return
		}

		// handle complete message
		if m.Attachments != nil && len(m.Attachments) > 0 && len(m.Content) > 0 {
			app.handleCompleteEvent(m)
			return
		}
	}
}

func (app *Application) messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.ChannelID != app.Config.Midjourney.ChannelID {
		return
	}

	log.Printf("messageDelete id: %s", m.ID)
}

// key: promptWithNoParameters
func (app *Application) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.ChannelID != app.Config.Midjourney.ChannelID {
		return
	}

	// maybe describe message update
	if m.Author == nil {
		if len(m.Embeds) > 0 && len(m.Attachments) == 0 && m.Embeds[0].Type == "rich" {
			app.handleDescribeUpdateEvent(m)
		}

		return
	}

	if m.Author.Username == "Midjourney Bot" {
		log.Printf("%s, messageUpdate: %s", m.ID, m.Content)

		if len(m.Attachments) > 0 && len(m.Content) > 0 {
			app.handleRateEvent(m)
			return
		}
	}
}

func webhookCallback(metaData *store.MetaData) {
	webhook := metaData.Webhook
	if webhook == "" {
		return
	}

	imageUrl, err := metaData.GetImageURL()
	if err != nil {
		log.Printf("Call metaData.GetImageURL failed, err: %+v", err)
		return
	}

	webhookReq := &wb.WebhookRequest{
		TaskID:       metaData.ID,
		Prompt:       metaData.Prompt,
		Type:         metaData.Type,
		Mode:         metaData.Mode,
		Status:       metaData.Status,
		ImageURL:     imageUrl,
		StartTime:    metaData.StartTime,
		CompleteTime: metaData.CompleteTime,
	}

	b, _ := json.Marshal(webhookReq)
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Call http.Post failed, %s, err: %+v", webhook, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Call http.Post failed, %s, status code: %d", webhook, resp.StatusCode)
		return
	}
}
