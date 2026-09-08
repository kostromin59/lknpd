package lknpd

import "time"

type LoginRequest struct {
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type LoginResponse struct {
	Token                string          `json:"token"`
	TokenExpireIn        time.Time       `json:"tokenExpireIn"`
	RefreshToken         string          `json:"refreshToken"`
	RefreshTokenExpireIn time.Time       `json:"refreshTokenExpireIn"`
	Profile              ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	INN string `json:"inn"`
}

type RefreshTokenRequest struct {
	DeviceInfo   DeviceInfo `json:"deviceInfo"`
	RefreshToken string     `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token         string    `json:"token"`
	TokenExpireIn time.Time `json:"tokenExpireIn"`
	RefreshToken  string    `json:"refreshToken"`
}
