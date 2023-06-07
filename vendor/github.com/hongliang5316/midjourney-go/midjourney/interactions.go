package midjourney

type InteractionsRequest struct {
	Type          int            `json:"type"`
	ApplicationID string         `json:"application_id"`
	MessageFlags  *int           `json:"message_flags,omitempty"`
	MessageID     *string        `json:"message_id,omitempty"`
	GuildID       string         `json:"guild_id"`
	ChannelID     string         `json:"channel_id"`
	SessionID     string         `json:"session_id"`
	Data          map[string]any `json:"data"`
}
