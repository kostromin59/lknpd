package lknpd

import (
	"encoding/base64"
	"encoding/json"
	"math/rand/v2"
	"strings"
)

const (
	deviceIDLength  = 21
	deviceIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// DeviceInfo is used to identify your device in lknpd.
type DeviceInfo struct {
	AppVersion     string      `json:"appVersion"`
	SourceDeviceID string      `json:"sourceDeviceId"`
	SourceType     string      `json:"sourceType"`
	MetaDetails    MetaDetails `json:"metaDetails"`
}

type MetaDetails struct {
	UserAgent string `json:"userAgent"`
}

func GenerateDeviceID() string {
	b := make([]byte, deviceIDLength)
	for i := range b {
		b[i] = deviceIDCharset[rand.N(len(deviceIDCharset))]
	}

	return string(b)
}

type tokenData struct {
	Sub any `json:"sub"`
}

type tokenDataSub struct {
	DeviceID       string `json:"deviceId"`
	RefreshContext struct {
		DeviceID string `json:"deviceId"`
	} `json:"refreshContext"`
}

func getDeviceIDFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		return ""
	}

	encodedData := parts[1]
	b, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return ""
	}

	var data tokenData
	_ = json.Unmarshal(b, &data)

	var subData tokenDataSub
	_ = json.Unmarshal([]byte(data.Sub.(string)), &subData)

	deviceID := subData.DeviceID
	if subData.RefreshContext.DeviceID != "" {
		deviceID = subData.RefreshContext.DeviceID
	}

	return deviceID
}
