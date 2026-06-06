package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT payload used for authentication.
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Role     string    `json:"role"`
	DeviceID string    `json:"device_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTMaker creates and verifies JWT tokens.
type JWTMaker struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTMaker(secret string, accessDuration, refreshDuration time.Duration) *JWTMaker {
	return &JWTMaker{
		secret:          []byte(secret),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

// CreateAccessToken generates a short-lived access token.
func (m *JWTMaker) CreateAccessToken(userID, tenantID uuid.UUID, role, deviceID string) (string, error) {
	return m.createToken(userID, tenantID, role, deviceID, m.accessDuration)
}

// CreateRefreshToken generates a long-lived refresh token.
func (m *JWTMaker) CreateRefreshToken(userID, tenantID uuid.UUID, role, deviceID string) (string, error) {
	return m.createToken(userID, tenantID, role, deviceID, m.refreshDuration)
}

func (m *JWTMaker) createToken(userID, tenantID uuid.UUID, role, deviceID string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("jwt.CreateToken: %w", err)
	}
	return signed, nil
}

// VerifyToken validates the token string and returns the decoded claims.
func (m *JWTMaker) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt.VerifyToken: unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.VerifyToken: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("jwt.VerifyToken: invalid token claims")
	}
	return claims, nil
}
