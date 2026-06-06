package auth

// ── DTOs ──────────────────────────────────────────────────────────────────────

type RegisterDTO struct {
	FullName string
	Email    string
	Password string
	FarmName string
}

type LoginDTO struct {
	Email    string
	Password string
	DeviceID string
}

type AuthResult struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	UserID          string `json:"user_id"`
	TenantID        string `json:"tenant_id"`
	Role            string `json:"role"`
	LastSyncVersion int64  `json:"last_sync_version"`
}
