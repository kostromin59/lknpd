package lknpd

import "math/rand/v2"

const (
	deviceIDLength  = 21
	deviceIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type DeviceInfo struct {
	AppVersion     string      `json:"appVersion"`
	SourceDeviceID string      `json:"sourceDeviceId"`
	SourceType     string      `json:"sourceType"`
	MetaDetails    MetaDetails `json:"metaDetails"`
}

type MetaDetails struct {
	UserAgent string `json:"userAgent"`
}

func generateDeviceID() string {
	b := make([]byte, deviceIDLength)
	for i := range b {
		b[i] = deviceIDCharset[rand.N(len(deviceIDCharset))]
	}

	return string(b)
}
