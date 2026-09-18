package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
)

type HMACTokenProvider struct{ secret []byte }
type payload struct {
	Subject string      `json:"sub"`
	Role    domain.Role `json:"role"`
	Expires int64       `json:"exp"`
	ID      string      `json:"jti"`
}

func NewHMACTokenProvider(secret string) *HMACTokenProvider {
	return &HMACTokenProvider{[]byte(secret)}
}
func (p *HMACTokenProvider) Generate(c ports.TokenClaims) (string, error) {
	b, err := json.Marshal(payload{c.Subject, c.Role, c.ExpiresAt.Unix(), c.TokenID})
	if err != nil {
		return "", err
	}
	body := base64.RawURLEncoding.EncodeToString(b)
	return body + "." + p.sign(body), nil
}
func (p *HMACTokenProvider) Parse(token string) (ports.TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(p.sign(parts[0])), []byte(parts[1])) {
		return ports.TokenClaims{}, domain.ErrInvalidToken
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return ports.TokenClaims{}, domain.ErrInvalidToken
	}
	var v payload
	if json.Unmarshal(b, &v) != nil || v.Subject == "" || time.Now().Unix() >= v.Expires {
		return ports.TokenClaims{}, domain.ErrInvalidToken
	}
	return ports.TokenClaims{Subject: v.Subject, Role: v.Role, ExpiresAt: time.Unix(v.Expires, 0), TokenID: v.ID}, nil
}
func (p *HMACTokenProvider) sign(v string) string {
	h := hmac.New(sha256.New, p.secret)
	_, _ = h.Write([]byte(v))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
