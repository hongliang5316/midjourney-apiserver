package store

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"github.com/redis/go-redis/v9"
)

type Error struct {
	Code api.Codes
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

type Store struct {
	*redis.Client
}

type Config struct {
	Redis Redis
}

func NewStore(config *Config) *Store {
	return &Store{
		Client: redis.NewClient(&redis.Options{
			Addr:     config.Redis.Address,
			Password: config.Redis.Password,
			DB:       0,
		}),
	}
}

type CompleteInfo struct {
	ID          string
	Attachments []map[string]any
}

func (s *Store) SaveWithComplete(ctx context.Context, completeMessageID, prompt, mode, attachments string, webhookFunc func(*MetaData)) error {
	key := GetKey(prompt)
	id, err := s.Get(ctx, key).Result()
	if err != nil {
		// not found, maybe expired or exception
		// drop it
		if err == redis.Nil {
			log.Printf("prompt: %s, attachments: %s complete, but no meta found", prompt, attachments)
			return nil
		}

		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Set failed, err: %+v", err),
		}
	}

	exist, err := s.Exists(ctx, id).Result()
	if err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Exists failed, err: %+v", err),
		}
	}

	if exist == 0 { // ignore this task
		return nil
	}

	if err := s.HSet(ctx, id, map[string]any{
		"complete_message_id": completeMessageID,
		"complete_time":       fmt.Sprint(time.Now().Unix()),
		"status":              string(StatusComplete),
		"mode":                mode,
		"process_rate":        "100%",
		"attachments":         attachments,
	}).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.HSet failed, err: %+v", err),
		}
	}

	metaData, err := s.GetMetaData(ctx, id)
	if err != nil {
		return err
	}

	if metaData == nil { // seems impossible
		return nil
	}

	webhookFunc(metaData)

	return nil
}

func GetKey(prompt string) string {
	index := strings.Index(prompt, "--")
	if index == -1 {
		return strings.TrimSpace(prompt)
	}

	return strings.TrimSpace(prompt[0:index])
}

func (s *Store) GetID(ctx context.Context, prompt string) (string, error) {
	key := GetKey(prompt)
	id, err := s.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}

		return "", Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Get failed, err: %+v", err),
		}
	}

	return id, nil
}

func (s *Store) CheckPrompt(ctx context.Context, prompt string) error {
	key := GetKey(prompt)
	id, err := s.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Set failed, err: %+v", err),
		}
	}

	val, err := s.HGet(ctx, id, "status").Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.HGet failed, err: %+v", err),
		}
	}

	if val != string(StatusComplete) && val != string(StatusJobQueued) {
		return Error{
			Code: api.Codes_CODES_INVALID_PARAMETER_ERROR,
			Msg:  fmt.Sprintf("The same prompt is being processed, please try again later."),
		}
	}

	return nil
}

func (s *Store) SaveWebhook(ctx context.Context, id, webhook string) error {
	exist, err := s.Exists(ctx, id).Result()
	if err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Exists failed, err: %+v", err),
		}
	}

	if exist == 0 { // ignore this task
		return nil
	}

	if err := s.HSet(ctx, id, map[string]any{
		"webhook": webhook,
	}).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.HSet failed, err: %+v", err),
		}
	}

	return nil
}

func (s *Store) SaveMeta(
	ctx context.Context,
	id,
	prompt string,
	status Status,
	typ Type,
	start_time int64,
) error {
	key := GetKey(prompt)
	if err := s.Set(ctx, key, id, Expired).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Set failed, err: %+v", err),
		}
	}

	if err := s.HSet(ctx, id, map[string]any{
		"id":         id,
		"status":     string(status),
		"type":       string(typ),
		"prompt":     prompt,
		"start_time": fmt.Sprint(start_time),
	}).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.HSet failed, err: %+v", err),
		}
	}

	if err := s.Expire(ctx, id, Expired).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Expire failed, err: %+v", err),
		}
	}

	return nil
}

func (s *Store) UpdateProcessRate(ctx context.Context, id, processRate string) error {
	exist, err := s.Exists(ctx, id).Result()
	if err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.Exists failed, err: %+v", err),
		}
	}

	if exist == 0 { // ignore this task
		return nil
	}

	if err := s.HSet(ctx, id, map[string]any{
		"status":       string(StatusProcessing),
		"process_rate": processRate,
	}).Err(); err != nil {
		return Error{
			Code: api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:  fmt.Sprintf("Call redis.HSet failed, err: %+v", err),
		}
	}

	return nil
}
