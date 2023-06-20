package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"github.com/hongliang5316/midjourney-go/midjourney"
)

var (
	Mutex          = new(sync.Mutex)
	DescribeInfoCh = make(chan discordgo.MessageEmbed, 1)
)

/*
flow:
1. create mesasge id: 1
2. update message id: 1
*/
func (s *Service) Describe(ctx context.Context, in *api.DescribeRequest) (*api.DescribeResponse, error) {
	if in.RequestId == "" {
		in.RequestId = uuid.NewString()
	}

	if ok := Mutex.TryLock(); !ok {
		return &api.DescribeResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_CONCURRENCY_LIMITED,
			Msg:       "Concurrency is limited",
		}, nil
	}

	defer Mutex.Unlock()

	clearCh()

	if err := s.MJClient.Describe(ctx, &midjourney.DescribeRequest{
		GuildID:   s.Config.Midjourney.GuildID,
		ChannelID: s.Config.Midjourney.ChannelID,
		ImageURL:  in.ImageUrl,
	}); err != nil {
		return &api.DescribeResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:       fmt.Sprint(err),
		}, nil
	}

	select {
	case <-time.After(30 * time.Second):
		return &api.DescribeResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_PROCESSING_TIMEOUT,
			Msg:       "timeout",
		}, nil
	case embedMsg := <-DescribeInfoCh:
		return &api.DescribeResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SUCCESS,
			Msg:       "success",
			Data: &api.DescribeResponseData{
				Prompts: getPrompts(embedMsg.Description),
			},
		}, nil
	}
}

func getPrompts(description string) []string {
	return strings.Split(description, "\n\n")
}

func clearCh() {
	for {
		select {
		case <-DescribeInfoCh:
		default:
			return
		}
	}
}
