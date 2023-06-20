package midjourney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type DescribeRequest struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	ImageURL  string `json:"image_url"`

	ext      string `json:"-"`
	filename string `json:"-"`
}

func downloadImage(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Call http.Get failed, err: %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Call ioutil.ReadAll failed, err: %+v", err)
	}

	return body, nil
}

func (c *Client) Describe(ctx context.Context, describeReq *DescribeRequest) error {
	if err := describeReq.init(); err != nil {
		return err
	}

	ext := strings.ToLower(describeReq.ext)
	if ext != "png" && ext != "jpg" {
		return fmt.Errorf("The image_url extension was only jpg and png formats are allowed currently")
	}

	image, err := downloadImage(ctx, describeReq.ImageURL)
	if err != nil {
		return err
	}

	attachmentsAndUploadResp, err := c.AttachmentsAndUpload(ctx, &AttachmentsAndUploadRequest{
		AttachmentsRequest: &AttachmentsRequest{
			ChannelID: describeReq.ChannelID,
			Files: []File{
				{
					ID:       "0",
					Filename: describeReq.filename,
					FileSize: int64(len(image)),
				},
			},
		},
		Image: image,
	})
	if err != nil {
		return fmt.Errorf("Call c.AttachmentsAndUpload failed, err: %w", err)
	}

	interactionsReq := &InteractionsRequest{
		Type:          2,
		ApplicationID: ApplicationID,
		GuildID:       describeReq.GuildID,
		ChannelID:     describeReq.ChannelID,
		SessionID:     SessionID,
		Data: map[string]any{
			"version": "1118961510123847774",
			"id":      "1092492867185950852",
			"name":    "describe",
			"type":    "1",
			"options": []map[string]any{
				{
					"type":  11,
					"name":  "image",
					"value": 0,
				},
			},
			"application_command": map[string]any{
				"id":                         "1092492867185950852",
				"application_id":             ApplicationID,
				"version":                    "1118961510123847774",
				"default_member_permissions": nil,
				"type":                       1,
				"nsfw":                       false,
				"name":                       "describe",
				"description":                "Writes a prompt based on your image.",
				"dm_permission":              true,
				"contexts":                   nil,
				"options": []map[string]any{
					{
						"type":        11,
						"name":        "image",
						"description": "The image to describe",
						"required":    true,
					},
				},
			},
			"attachments": []*AttachmentsAndUploadResponse{
				attachmentsAndUploadResp,
			},
		},
	}

	b, _ := json.Marshal(interactionsReq)

	url_ := "https://discord.com/api/v9/interactions"
	req, err := http.NewRequest("POST", url_, bytes.NewReader(b))
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

func (r *DescribeRequest) init() error {
	u, err := url.Parse(r.ImageURL)
	if err != nil {
		return fmt.Errorf("Call url.Parse failed, err: %w", err)
	}

	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return fmt.Errorf("Couldn't find a period to indicate a file extension")
	}

	r.ext = u.Path[pos+1 : len(u.Path)]
	r.filename = path.Base(u.Path)
	return nil
}
