package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
)

type WechatService struct {
	appID     string
	appSecret string
}

func NewWechatService() *WechatService {
	cfg := conf.Global()
	return &WechatService{
		appID:     cfg.WECHAT.AppID,
		appSecret: cfg.WECHAT.AppSecret,
	}
}

type Code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (s *WechatService) Code2Session(ctx context.Context, code string) (*Code2SessionResp, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperrors.ErrWechatAPIFailed.Wrap(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, apperrors.ErrWechatAPIFailed.Wrap(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperrors.ErrWechatAPIFailed.Wrap(err)
	}

	var result Code2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, apperrors.ErrWechatAPIFailed.Wrap(err)
	}

	if result.ErrCode != 0 {
		return nil, apperrors.ErrWechatCodeInvalid
	}

	return &result, nil
}
