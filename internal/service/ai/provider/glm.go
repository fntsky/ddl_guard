package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/ai/shared"
)

const providerName = "glm"

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
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "build_image_url", err)
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
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "marshal_request", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return schema.CreateDraftResp{}, apperrors.AIRequestFailed(providerName, 0, "", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return schema.CreateDraftResp{}, apperrors.AIRequestFailed(providerName, 0, "", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "read_response", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return schema.CreateDraftResp{}, apperrors.AIRequestFailed(providerName, resp.StatusCode, string(body), nil)
	}

	chatResp := struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}{}
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "parse_response", err)
	}
	if len(chatResp.Choices) == 0 {
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "no_choices", nil)
	}

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if content == "" {
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "empty_content", nil)
	}

	result, err := shared.ParseDraftFromModelJSON(content)
	if err != nil {
		return schema.CreateDraftResp{}, apperrors.AIResponseInvalid(providerName, "parse_draft", err)
	}

	return result, nil
}