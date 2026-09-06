package lknpd

import "math/rand/v2"

const (
	deviceIDLength  = 21
	deviceIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateDeviceID() string {
	b := make([]byte, deviceIDLength)
	for i := range b {
		b[i] = deviceIDCharset[rand.N(len(deviceIDCharset))]
	}

	return string(b)
}
