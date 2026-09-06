package lknpd

import "math/rand/v2"

const (
	deviceIDLength  = 21
	deviceIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type DeviceInfo struct {
	AppVersion     string      `json:"appVersion,omitempty"`
	SourceDeviceID string      `json:"sourceDeviceId,omitempty"`
	SourceType     string      `json:"sourceType,omitempty"`
	MetaDetails    MetaDetails `json:"metaDetails,omitempty"`
}

type MetaDetails struct {
	UserAgent string `json:"userAgent,omitempty"`
}

func generateDeviceID() string {
	b := make([]byte, deviceIDLength)
	for i := range b {
		b[i] = deviceIDCharset[rand.N(len(deviceIDCharset))]
	}

	return string(b)
}
