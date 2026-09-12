package lknpd_test

import (
	"testing"

	"github.com/kostromin59/lknpd"
)

func TestConstructor(t *testing.T) {
	t.Run("initial generated device info", func(t *testing.T) {
		c := lknpd.New()
		deviceInfo := c.DeviceInfo()

		if deviceInfo.AppVersion == "" {
			t.Error("expected non-empty deviceInfo.AppVersion")
		}

		if deviceInfo.SourceDeviceID == "" {
			t.Error("expected non-empty deviceInfo.SourceDeviceID")
		}

		if deviceInfo.SourceType == "" {
			t.Error("expected non-empty deviceInfo.SourceType")
		}

		if deviceInfo.MetaDetails.UserAgent == "" {
			t.Error("expected non-empty deviceInfo.MetaDetails.UserAgent")
		}
	})

	t.Run("custom options", func(t *testing.T) {
		expectedAppVersion := "testAppVersion"
		expectedDeviceID := "testDeviceID"
		expectedSourceType := "testSourceType"
		expectedUserAgent := "testUserAgent"

		c := lknpd.New(
			lknpd.WithAppVersion(expectedAppVersion),
			lknpd.WithDeviceID(expectedDeviceID),
			lknpd.WithSourceType(expectedSourceType),
			lknpd.WithUserAgent(expectedUserAgent),
		)

		deviceInfo := c.DeviceInfo()

		if deviceInfo.AppVersion != expectedAppVersion {
			t.Errorf("expected deviceInfo.AppVersion %q but got %q", expectedAppVersion, deviceInfo.AppVersion)
		}

		if deviceInfo.SourceDeviceID != expectedDeviceID {
			t.Errorf("expected deviceInfo.DeviceID %q but got %q", expectedDeviceID, deviceInfo.SourceDeviceID)
		}

		if deviceInfo.SourceType != expectedSourceType {
			t.Errorf("expected deviceInfo.SourceType %q but got %q", expectedSourceType, deviceInfo.SourceType)
		}

		if deviceInfo.MetaDetails.UserAgent != expectedUserAgent {
			t.Errorf("expected deviceInfo.MetaDetails.UserAgent %q but got %q", expectedUserAgent, deviceInfo.MetaDetails.UserAgent)
		}
	})

	t.Run("custom tokens", func(t *testing.T) {
		expectedToken := "testToken"
		expectedRefreshToken := "testRefreshToken"

		c := lknpd.New(
			lknpd.WithToken(expectedToken),
			lknpd.WithRefreshToken(expectedRefreshToken),
		)
		deviceInfo := c.DeviceInfo()

		if c.Token() != expectedToken {
			t.Errorf("expected c.Token %q but got %q", expectedToken, c.Token())
		}

		if c.RefreshToken() != expectedRefreshToken {
			t.Errorf("expected c.RefreshToken %q but got %q", expectedRefreshToken, c.RefreshToken())
		}

		if deviceInfo.SourceDeviceID == "" {
			t.Error("expected non-empty deviceInfo.SourceDeviceID")
		}
	})

	t.Run("deviceID from token", func(t *testing.T) {
		expectedToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJkZXZpY2VJZFwiOlwidGVzdFRva2VuRGV2aWNlSURcIn0iLCJleHAiOjE3ODkyMDAwMDB9.o7_wMC2Y6sIdxCQ1KMKXMNd1n85renPib2gKxRys2queUXx6qJU8BK9KhLz264LlGZzgbjuPSIzUwusY1fH5Og"
		expectedDeviceID := "testTokenDeviceID"

		c := lknpd.New(
			lknpd.WithToken(expectedToken),
		)
		deviceInfo := c.DeviceInfo()

		if c.Token() != expectedToken {
			t.Errorf("expected c.Token %q but got %q", expectedToken, c.Token())
		}

		if deviceInfo.SourceDeviceID != expectedDeviceID {
			t.Errorf("expected deviceInfo.DeviceID %q but got %q", expectedDeviceID, deviceInfo.SourceDeviceID)
		}
	})

	t.Run("deviceID from refresh token", func(t *testing.T) {
		expectedToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJkZXZpY2VJZFwiOlwidGVzdFRva2VuRGV2aWNlSURcIn0iLCJleHAiOjE3ODkyMDAwMDB9.o7_wMC2Y6sIdxCQ1KMKXMNd1n85renPib2gKxRys2queUXx6qJU8BK9KhLz264LlGZzgbjuPSIzUwusY1fH5Og"
		expectedRefreshToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJyZWZyZXNoQ29udGV4dFwiOiB7XCJkZXZpY2VJZFwiOlwidGVzdFJlZnJlc2hUb2tlbkRldmljZUlEXCJ9fSIsImV4cCI6MTc4OTIwMDAwMH0.a30vwjcAr2pHwh6EPQoPV_0Vhehpj2iGLLCq3UnkiyYoxdeMHhkkBK1moJZ5LrpejeDQtL0BlSe9DrdgRSNdmA"
		expectedDeviceID := "testRefreshTokenDeviceID"

		c := lknpd.New(
			lknpd.WithToken(expectedToken),
			lknpd.WithRefreshToken(expectedRefreshToken),
		)
		deviceInfo := c.DeviceInfo()

		if c.RefreshToken() != expectedRefreshToken {
			t.Errorf("expected c.RefreshToken %q but got %q", expectedRefreshToken, c.RefreshToken())
		}

		if deviceInfo.SourceDeviceID != expectedDeviceID {
			t.Errorf("expected deviceInfo.DeviceID %q but got %q", expectedDeviceID, deviceInfo.SourceDeviceID)
		}
	})
}
