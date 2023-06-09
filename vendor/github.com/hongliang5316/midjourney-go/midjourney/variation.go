package midjourney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type VariationRequest struct {
	Index       int32  `json:"index"`
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	MessageID   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (c *Client) Variation(ctx context.Context, variationReq *VariationRequest) error {
	flags := 0
	interactionsReq := &InteractionsRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       variationReq.GuildID,
		ChannelID:     variationReq.ChannelID,
		MessageFlags:  &flags,
		MessageID:     &variationReq.MessageID,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      fmt.Sprintf("MJ::JOB::variation::%d::%s", variationReq.Index, variationReq.MessageHash),
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
