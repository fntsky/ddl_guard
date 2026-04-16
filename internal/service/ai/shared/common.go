package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/schema"
)

const visionPrompt = "You are a DDL assistant. Read the image and output STRICT JSON only with keys: title, description, deadline, early_remind. deadline must be RFC3339 format with timezone (example: 2026-03-24T14:30:00+08:00). early_remind must be an integer number of minutes. Do not output markdown."

type draftVisionResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	EarlyRemind int    `json:"early_remind"`
}

func VisionPrompt() string {
	return visionPrompt
}

func BuildImageDataURL(imageData []byte) (string, error) {
	if len(imageData) == 0 {
		return "", fmt.Errorf("empty image data")
	}

	mime := http.DetectContentType(imageData)
	if !strings.HasPrefix(mime, "image/") {
		mime = "image/png"
	}

	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(imageData)), nil
}

func ParseDraftFromModelJSON(raw string) (schema.CreateDraftResp, error) {
	out := draftVisionResult{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &out); err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("parse model json failed: %w", err)
	}

	deadline, err := time.Parse(time.RFC3339, strings.TrimSpace(out.Deadline))
	if err != nil {
		return schema.CreateDraftResp{}, fmt.Errorf("invalid deadline format: %w", err)
	}

	return schema.CreateDraftResp{
		Title:       strings.TrimSpace(out.Title),
		Description: strings.TrimSpace(out.Description),
		Deadline:    deadline,
		EarlyRemind: out.EarlyRemind,
	}, nil
}
