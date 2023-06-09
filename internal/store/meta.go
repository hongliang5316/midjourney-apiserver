package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
)

type MetaData struct {
	ID                string `redis:"id"`
	Prompt            string `redis:"prompt"`
	Type              Type   `redis:"type"`
	Status            Status `redis:"status"`
	ProcessRate       string `redis:"process_rate"`
	Attachments       string `redis:"attachments"`
	StartTime         int64  `redis:"start_time"`
	CompleteTime      int64  `redis:"complete_time"`
	Webhook           string `redis:"webhook"`
	CompleteMessageID string `redis:"complete_message_id"`
	Mode              string `redis:"mode"`
}

func (md *MetaData) GetImageURL() (string, error) {
	as := []discordgo.MessageAttachment{}
	if err := json.Unmarshal([]byte(md.Attachments), &as); err != nil {
		return "", Error{
			Code: api.Codes_CODES_SERVER_ERROR,
			Msg:  fmt.Sprintf("Call json.Unmarshal failed, err: %+v", err),
		}
	}

	return as[0].URL, nil
}

func (s *Store) GetMetaData(ctx context.Context, id string) (*MetaData, error) {
	res := s.HGetAll(ctx, id)
	if res.Err() != nil {
		return nil, Error{
			Code: api.Codes_CODES_SERVER_ERROR,
			Msg:  fmt.Sprintf("Call s.HGetAll failed, err: %+v", res.Err()),
		}
	}

	// empty array
	if len(res.Val()) == 0 {
		return nil, nil
	}

	md := new(MetaData)
	if err := res.Scan(md); err != nil {
		return nil, Error{
			Code: api.Codes_CODES_SERVER_ERROR,
			Msg:  fmt.Sprintf("Call res.Scan failed, err: %+v", err),
		}
	}

	return md, nil
}
