package midjourney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ImagineRequest struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	Prompt    string `json:"prompt"`
}

func (c *Client) Imagine(ctx context.Context, imgReq *ImagineRequest) error {
	interactionsReq := &InteractionsRequest{
		Type:          2,
		ApplicationID: ApplicationID,
		GuildID:       imgReq.GuildID,
		ChannelID:     imgReq.ChannelID,
		SessionID:     SessionID,
		Data: map[string]any{
			"version": "1077969938624553050",
			"id":      "938956540159881230",
			"name":    "imagine",
			"type":    "1",
			"options": []map[string]any{
				{
					"type":  3,
					"name":  "prompt",
					"value": imgReq.Prompt,
				},
			},
			"application_command": map[string]any{
				"id":                         "938956540159881230",
				"application_id":             ApplicationID,
				"version":                    "1077969938624553050",
				"default_permission":         true,
				"default_member_permissions": nil,
				"type":                       1,
				"nsfw":                       false,
				"name":                       "imagine",
				"description":                "Create images with Midjourney",
				"dm_permission":              true,
				"options": []map[string]any{
					{
						"type":        3,
						"name":        "prompt",
						"description": "The prompt to imagine",
						"required":    true,
					},
				},
				"attachments": []any{},
			},
		},
	}

	b, _ := json.Marshal(interactionsReq)

	url := "https://discord.com/api/v9/interactions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Call http.NewRequest failed, err: %w", err)
	}

	req.Header.Set("Authorization", c.Config.UserToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Call c.Do failed, err: %w", err)
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("Call checkResponse failed, err: %w", err)
	}

	return nil
}
