package midjourney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AttachmentsResponse struct {
	Attachments []Attachment
}

type Attachment struct {
	ID             int64  `json:"id"`
	UploadURL      string `json:"upload_url"`
	UploadFilename string `json:"upload_filename"`
}

type AttachmentsRequest struct {
	ChannelID string `json:"channel_id"`
	Files     []File `json:"files"`
}

type File struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
}

type AttachmentsAndUploadRequest struct {
	*AttachmentsRequest
	Image []byte
}

type AttachmentsAndUploadResponse struct {
	ID               string `json:"id"`
	Filename         string `json:"filename"`
	UploadedFilename string `json:"uploaded_filename"`
}

func (c *Client) AttachmentsAndUpload(ctx context.Context, attachmentsAndUploadReq *AttachmentsAndUploadRequest) (*AttachmentsAndUploadResponse, error) {
	attachmentsResp, err := c.Attachments(ctx, attachmentsAndUploadReq.AttachmentsRequest)
	if err != nil {
		return nil, fmt.Errorf("Call c.Attachments failed, err: %w", err)
	}

	url := attachmentsResp.Attachments[0].UploadURL
	req, err := http.NewRequest("PUT", url, bytes.NewReader(attachmentsAndUploadReq.Image))
	if err != nil {
		return nil, fmt.Errorf("Call http.NewRequest failed, err: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Call c.Do failed, err: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return nil, fmt.Errorf("Bad http status code %d", resp.StatusCode)
	}

	return &AttachmentsAndUploadResponse{
		ID:               attachmentsAndUploadReq.Files[0].ID,
		Filename:         attachmentsAndUploadReq.Files[0].Filename,
		UploadedFilename: attachmentsResp.Attachments[0].UploadFilename,
	}, nil
}

func (c *Client) Attachments(ctx context.Context, attachmentsReq *AttachmentsRequest) (*AttachmentsResponse, error) {
	if len(attachmentsReq.Files) != 1 {
		return nil, fmt.Errorf("only support one image to upload currently")
	}

	b, _ := json.Marshal(attachmentsReq)

	url := fmt.Sprintf("https://discord.com/api/v9/channels/%s/attachments", attachmentsReq.ChannelID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("Call http.NewRequest failed, err: %w", err)
	}

	req.Header.Set("Authorization", c.Config.UserToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Call c.Do failed, err: %w", err)
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("Call checkResponse failed, err: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Call ioutil.ReadAll failed, err: %w", err)
	}

	attachmentsResp := new(AttachmentsResponse)
	if err := json.Unmarshal(body, attachmentsResp); err != nil {
		return nil, fmt.Errorf("Call json.Unmarshal failed, body: %s, err: %w", string(body), err)
	}

	return attachmentsResp, nil
}
