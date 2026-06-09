package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Session TTL sourced from config (fallback: jwt.expire).
const sessionExpireDefault = 14 // days

// Redis key templates.
const (
	redisSessionPrefix    = "session:"          // session:{sid} -> JSON payload
	redisUserSessionIdx  = "user_session:"      // user_session:{userId} -> sid
	sessionCookieName    = "suask_sid"
)

// SessionPayload stored in Redis for each active session.
type SessionPayload struct {
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"` // unix milli
}

// GenerateSID creates a cryptographically random 64-char hex string.
func GenerateSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("session: generate sid failed: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// getExpireSeconds reads the session expire config in days and converts to seconds.
// Falls back to jwt.expire if session.expire is not set.
func getExpireSeconds() int64 {
	d := g.Cfg().MustGet(context.TODO(), "session.expire", g.Cfg().MustGet(context.TODO(), "jwt.expire", sessionExpireDefault)).Int64()
	return d * 24 * 3600
}

// CreateSession writes a new session to Redis and indexes it by userId.
// Returns the SID.
func CreateSession(ctx context.Context, userID int, role string) (string, error) {
	sid, err := GenerateSID()
	if err != nil {
		return "", err
	}

	now := time.Now().UnixMilli()
	payload := SessionPayload{
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("session: marshal payload failed: %w", err)
	}

	ex := getExpireSeconds()
	sessionKey := redisSessionPrefix + sid

	// Write session data.
	if err = g.Redis().SetEX(ctx, sessionKey, string(data), ex); err != nil {
		return "", fmt.Errorf("session: setex %s failed: %w", sessionKey, err)
	}

	// Index: userId -> sid (also sets TTL so stale index keys expire).
	idxKey := redisUserSessionIdx + strconv.Itoa(userID)
	if err = g.Redis().SetEX(ctx, idxKey, sid, ex); err != nil {
		return "", fmt.Errorf("session: setex %s failed: %w", idxKey, err)
	}

	return sid, nil
}

// GetSession loads a session from Redis by SID. Returns nil if not found.
func GetSession(ctx context.Context, sid string) (*SessionPayload, error) {
	key := redisSessionPrefix + sid
	v, err := g.Redis().Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("session: get %s failed: %w", key, err)
	}
	if v.IsNil() || v.String() == "" {
		return nil, nil
	}

	var p SessionPayload
	if err := json.Unmarshal([]byte(v.String()), &p); err != nil {
		return nil, fmt.Errorf("session: unmarshal %s failed: %w", key, err)
	}
	return &p, nil
}

// DeleteSession removes a session from Redis and its user index.
func DeleteSession(ctx context.Context, sid string) error {
	key := redisSessionPrefix + sid
	// Read userID from session before deleting.
	p, err := GetSession(ctx, sid)
	if err != nil {
		return err
	}
	if p != nil {
		idxKey := redisUserSessionIdx + strconv.Itoa(p.UserID)
		_, _ = g.Redis().Del(ctx, idxKey)
	}
	_, err = g.Redis().Del(ctx, key)
	return err
}

// DeleteUserSessions removes all sessions for a given user.
func DeleteUserSessions(ctx context.Context, userID int) error {
	idxKey := redisUserSessionIdx + strconv.Itoa(userID)
	v, err := g.Redis().Get(ctx, idxKey)
	if err != nil {
		return err
	}
	if !v.IsNil() && v.String() != "" {
		_ = DeleteSession(ctx, v.String())
	}
	_, _ = g.Redis().Del(ctx, idxKey)
	return nil
}

// RefreshSessionTTL extends the TTL of a session if it's below the threshold.
// Returns the new TTL in seconds, or 0 if no refresh was needed.
func RefreshSessionTTL(ctx context.Context, sid string) (int64, error) {
	key := redisSessionPrefix + sid
	maxTTL := getExpireSeconds()
	threshold := maxTTL / 2

	ttl, err := g.Redis().TTL(ctx, key)
	if err != nil {
		return 0, err
	}
	if ttl <= 0 || ttl > threshold {
		return ttl, nil
	}
	if _, expireErr := g.Redis().Expire(ctx, key, maxTTL); expireErr != nil {
		return 0, expireErr
	}
	// Also refresh the index TTL.
	p, err := GetSession(ctx, sid)
	if err == nil && p != nil {
		idxKey := redisUserSessionIdx + strconv.Itoa(p.UserID)
		_, _ = g.Redis().Expire(ctx, idxKey, maxTTL)
	}
	return maxTTL, nil
}

// CookieName returns the session cookie name.
func CookieName() string {
	return sessionCookieName
}

// CookieConfig returns the cookie name and max-age (seconds) for setting headers.
func CookieConfig() (name string, maxAge int) {
	return sessionCookieName, int(getExpireSeconds())
}
