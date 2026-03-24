package ai

import (
	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/ai/provider"
)

type AIProvider interface {
	AnalyzeImage(imageData []byte) (schema.CreateDraftResp, error)
}

func NewAIProvider() AIProvider {
	global := conf.Global()
	if global == nil {
		return nil
	}
	config := global.VISUAL_AI

	switch config.Provider {
	case "glm":
		return provider.NewGLMProvider(config)
	default:
		return nil
	}
}
