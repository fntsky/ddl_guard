# Spec: WeChat Mini Program Login/Register

## Objective

Add a WeChat Mini Program login/register endpoint to the existing auth system. When a user logs in via WeChat for the first time, a new account is automatically created (register-on-login). Subsequent logins return tokens for the existing account.

**User story:** As a WeChat Mini Program user, I can tap "WeChat Login" on the client, the client calls `wx.login()` to get a `js_code`, sends it to our backend, and I receive JWT tokens — no separate registration step needed.

**Flow:**
1. Client calls `wx.login()` → gets `js_code`
2. Client sends `POST /api/v1/users/sessions/wechat` with `{ "code": "js_code" }`
3. Backend calls WeChat `code2Session` API → gets `openid` + `session_key`
4. Backend looks up `UserAuth` by `AuthType=wechat, AuthIdentifier=openid`
   - If found: load user, issue tokens (login)
   - If not found: create new user + UserAuth record, issue tokens (register)
5. Return `{ uuid, username, access_token, refresh_token }`

## Tech Stack

- Go 1.x, Gin HTTP framework, XORM ORM, PostgreSQL
- Google Wire for DI
- WeChat code2Session API: `GET https://api.weixin.qq.com/sns/jscode2session`

## Commands

```bash
# Build
go build -o ddl_guard ./cmd/ddl_guard/

# Regenerate Wire (required after adding new providers)
go generate ./cmd/...

# Run
./ddl_guard run -d ./data

# DB migration
./ddl_guard upgrade -d ./data
```

## Project Structure (changes only)

```
internal/
  ├── base/conf/config.go          → Add WechatConfig struct
  ├── entity/user_auth_entity.go   → Add UserAuthTypeWechat constant
  ├── repo/user_auth/              → NEW: UserAuth repository
  │   └── user_auth_repo.go
  ├── service/wechat/              → NEW: WeChat code2Session client
  │   └── wechat_service.go
  ├── service/user/user_service.go → Add LoginByWechat method
  ├── controller/user_controller.go → Add LoginByWechat handler
  ├── schema/user_schema.go        → Add WeChat request/response DTOs
  ├── router/user_api_router.go    → Add WeChat route
  ├── errors/errors.go             → Add WeChat error codes
  ├── migrations/migrations.go     → Add migration v0.0.8 (no schema change needed, UserAuth table already exists)
configs/config.yaml                → Add WECHAT config section
```

## Code Style

Follow existing patterns exactly. Example — the new WeChat login mirrors the email login pattern:

```go
// schema/user_schema.go — request/response pattern
type LoginByWechatReq struct {
    Code string `json:"code" binding:"required"`
}

type LoginByWechatResp struct {
    UUID         string `json:"uuid"`
    Username     string `json:"username"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// service — method pattern
func (s *UserService) LoginByWechat(ctx context.Context, req *schema.LoginByWechatReq) (*schema.LoginByWechatResp, error) {
    // 1. Call WeChat API to exchange code for openid
    // 2. Look up UserAuth by (wechat, openid)
    // 3. If found → login; if not → create user + auth record
    // 4. Issue tokens
}
```

## Key Design Decisions

### 1. Use existing `UserAuth` table for WeChat identity
The `UserAuth` entity already exists with `AuthType`/`AuthIdentifier`/`AuthMeta` designed for multi-provider auth. We add `UserAuthTypeWechat = "wechat"` and store `openid` in `AuthIdentifier`, `session_key` in `AuthMeta`.

### 2. Auto-register on first login (no separate register endpoint)
WeChat users don't set passwords or usernames on first login. The `User` record is created with:
- `Username`: `"wx_" + first 8 chars of openid` (placeholder, user can update later)
- `Email`: nil
- `Phone`: nil
- `PasswordHash`: random bcrypt hash (unusable, prevents password-based login)

### 3. WeChat API client as a standalone service
A `WechatService` encapsulates the HTTP call to `code2Session`. It reads `appid`/`appsecret` from config. This keeps the WeChat API concern separate from user business logic.

### 4. No migration needed
The `UserAuth` table is already synced in `init_data.go`. No schema changes required.

### 5. Config via YAML + env override
Following existing pattern (like `JWT` and `EMAIL_OTP`), add a `WECHAT` section to config with `env` tags for secret override.

## Testing Strategy

- Manual testing via Swagger UI or curl
- Verify: new user creation on first login, existing user login on subsequent logins
- Verify: invalid `js_code` returns appropriate error
- Verify: WeChat API errors (invalid appid, etc.) are handled gracefully

## Boundaries

- **Always:** Follow existing code patterns, use `handler.HandleResponse`, use `AppError` for errors, add Swagger annotations
- **Ask first:** Adding new dependencies, changing existing entity schemas
- **Never:** Store `appsecret` in client-facing code, expose `session_key` to the client, modify JWT token structure

## Success Criteria

1. `POST /api/v1/users/sessions/wechat` with valid `code` returns JWT tokens
2. First-time WeChat user gets auto-registered with a `UserAuth` record
3. Returning WeChat user gets logged in without creating duplicate accounts
4. Invalid/expired `code` returns 401 with clear error message
5. WeChat config (appid, appsecret) is loaded from config file with env override
6. Wire DI is regenerated and compiles successfully
7. Swagger docs include the new endpoint

## Open Questions

- Should we support WeChat account binding to an existing email/phone account? (Deferred — not in scope for this iteration)
- Should the auto-generated username be updatable? (Yes, but the update endpoint is out of scope — the `username` field is already on the `User` entity)
