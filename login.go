package lknpd

import "time"

type LoginRequest struct {
	Username   string     `json:"username,omitempty"`
	Password   string     `json:"password,omitempty"`
	DeviceInfo DeviceInfo `json:"deviceInfo,omitempty"`
}

type LoginResponse struct {
	Token          string    `json:"token,omitempty"`
	TokenExpiresIn time.Time `json:"tokenExpiresIn,omitempty"`
	RefreshToken   string    `json:"refreshToken,omitempty"`
}

type RefreshTokenRequest struct {
	DeviceInfo   DeviceInfo `json:"deviceInfo,omitempty"`
	RefreshToken string     `json:"refreshToken,omitempty"`
}

type RefreshTokenResponse struct {
	Token          string    `json:"token,omitempty"`
	TokenExpiresIn time.Time `json:"tokenExpiresIn,omitempty"`
	RefreshToken   string    `json:"refreshToken,omitempty"`
}
