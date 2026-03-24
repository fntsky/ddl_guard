package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/ai/shared"
)

type GLMProvider struct {
	APIKey   string
	Endpoint string
	Model    string
	Client   *http.Client
}

func NewGLMProvider(config conf.VisualAIConfig) *GLMProvider {
	return &GLMProvider{
		APIKey:   config.APIKey,
		Endpoint: config.Endpoint,
		Model:    config.Model,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *GLMProvider) AnalyzeImage(imageData []byte) (schema.CreateDraftResp, error) {
	dataURL, err := shared.BuildImageDataURL(imageData)
	if err != nil {
		return schema.CreateDraftResp{}, err
	}

	model := strings.TrimSpace(p.Model)
	if model == "" {
		model = "glm-4.6v"
	}

	// Keep payload shape close to GLM official sample.
	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "image_url",
						"image_url": map[string]string{
							// Base64 data URL.
							"url": dataURL,
						},
					},
					{
						"type": "text",
						"text": shared.VisionPrompt(),
					},
				},
			},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("marshal glm request failed: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("create glm request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("glm request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("read glm response failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return schema.CreateDraftResp{}, fmt.Errorf("glm api failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	chatResp := struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}{}
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("unmarshal glm response failed: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return schema.CreateDraftResp{}, fmt.Errorf("glm response has no choices")
	}

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if content == "" {
		return schema.CreateDraftResp{}, fmt.Errorf("glm response content is empty")
	}

	return shared.ParseDraftFromModelJSON(content)
}
