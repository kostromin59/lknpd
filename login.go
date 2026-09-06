package lknpd

import "time"

type LoginRequest struct {
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type LoginResponse struct {
	Token          string    `json:"token"`
	TokenExpiresIn time.Time `json:"tokenExpiresIn"`
	RefreshToken   string    `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	DeviceInfo   DeviceInfo `json:"deviceInfo"`
	RefreshToken string     `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token          string    `json:"token"`
	TokenExpiresIn time.Time `json:"tokenExpiresIn"`
	RefreshToken   string    `json:"refreshToken"`
}
