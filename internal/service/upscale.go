package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"github.com/hongliang5316/midjourney-apiserver/pkg/store"
	"github.com/hongliang5316/midjourney-go/midjourney"
)

/*
flow:
1. create mesasge id: 3
2. update message id: 2
3. create message id: 4 -> contains attachments
4. delete message id: 3
*/
func (s *Service) Upscale(ctx context.Context, in *api.UpscaleRequest) (*api.UpscaleResponse, error) {
	if in.RequestId == "" {
		in.RequestId = uuid.NewString()
	}

	metaData, err := s.Store.GetMetaData(ctx, in.TaskId)
	if err != nil {
		e := err.(store.Error)
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      e.Code,
			Msg:       e.Msg,
		}, nil
	}

	if metaData == nil {
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_INVALID_PARAMETER_ERROR,
			Msg:       fmt.Sprintf("id: %s, not found", in.TaskId),
		}, nil
	}

	if metaData.Type != store.TypeImagine {
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_INVALID_PARAMETER_ERROR,
			Msg:       fmt.Sprintf("id: %s, the type is not `Imagine`, %s", in.TaskId, metaData.Type),
		}, nil
	}

	if metaData.Status != store.StatusComplete {
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_INVALID_PARAMETER_ERROR,
			Msg:       fmt.Sprintf("id: %s, the status is not `Complete`, %s", in.TaskId, metaData.Status),
		}, nil
	}

	url, err := metaData.GetImageURL()
	if err != nil {
		e := err.(store.Error)
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      e.Code,
			Msg:       e.Msg,
		}, nil
	}

	key := store.GetKey(metaData.Prompt)

	log.Printf("Upscale, key: %s, len: %d", key, len(key))

	KeyChan.Init(key)
	defer KeyChan.Del(key)

	if err := s.MJClient.Upscale(ctx, &midjourney.UpscaleRequest{
		Index:       in.Index,
		GuildID:     s.Config.Midjourney.GuildID,
		ChannelID:   s.Config.Midjourney.ChannelID,
		MessageID:   metaData.CompleteMessageID,
		MessageHash: midjourney.GetMessageHash(url),
	}); err != nil {
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:       fmt.Sprint(err),
		}, nil
	}

	select {
	case <-time.After(10 * time.Second):
		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:       "timeout",
		}, nil
	case msgInfo := <-KeyChan.Get(key):
		if msgInfo.Error != nil {
			code := api.Codes_CODES_SERVER_INTERNAL_ERROR

			switch msgInfo.Error.Title {
			case "Invalid parameter":
				code = api.Codes_CODES_INVALID_PARAMETER_ERROR
			}

			return &api.UpscaleResponse{
				RequestId: in.RequestId,
				Code:      code,
				Msg:       msgInfo.Error.Description,
			}, nil
		}

		if err := s.Store.SaveWebhook(ctx, msgInfo.ID, in.Webhook); err != nil {
			return &api.UpscaleResponse{
				RequestId: in.RequestId,
				Code:      api.Codes_CODES_SERVER_INTERNAL_ERROR,
				Msg:       fmt.Sprint(err),
			}, nil
		}

		return &api.UpscaleResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SUCCESS,
			Msg:       "success",
			Data: &api.UpscaleResponseData{
				TaskId:    msgInfo.ID,
				StartTime: msgInfo.StartTime,
			},
		}, nil
	}
}
