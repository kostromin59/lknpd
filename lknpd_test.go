package lknpd_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

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

func TestRequestWithAuth(t *testing.T) {
	expectedRefreshToken := "refreshToken"
	expectedToken := "token"

	expectedNewRefreshToken := "newRefreshToken"
	expectedNewToken := "newToken"
	expectedExpireIn := time.Now().Add(5 * time.Minute)

	oldRefreshToken := "oldRefreshToken"

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/token", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		var request lknpd.RefreshTokenRequest
		_ = json.NewDecoder(r.Body).Decode(&request)

		if request.RefreshToken == oldRefreshToken {
			_ = json.NewEncoder(w).Encode(lknpd.RefreshTokenResponse{
				Token:         "some",
				TokenExpireIn: time.Now().Add(-1 * time.Hour),
				RefreshToken:  expectedRefreshToken,
			})
			w.WriteHeader(http.StatusCreated)
			return
		}

		if request.RefreshToken != expectedRefreshToken {
			t.Errorf("expected refresh token %q but got %q", expectedRefreshToken, request.RefreshToken)
		}

		_ = json.NewEncoder(w).Encode(lknpd.RefreshTokenResponse{
			Token:         expectedNewToken,
			TokenExpireIn: expectedExpireIn,
			RefreshToken:  expectedNewRefreshToken,
		})
		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("GET /token", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+expectedToken {
			t.Errorf("expected Authorization Header %q but got %q", "Bearer "+expectedToken, r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+expectedNewToken {
			t.Errorf("expected Authorization Header %q but got %q", "Bearer "+expectedNewToken, r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
	})

	mux.HandleFunc("POST /error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(lknpd.ErrorResponse{
			Code:    "errorCode",
			Message: "someMsg",
		})
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	t.Run("successful request", func(t *testing.T) {
		c := lknpd.New(
			lknpd.WithBaseURL(testServer.URL),
			lknpd.WithToken(expectedToken),
		)

		_, err := c.RequestWithAuth[any](t.Context(), "/token", http.MethodGet, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("bad status code", func(t *testing.T) {
		c := lknpd.New(
			lknpd.WithBaseURL(testServer.URL),
			lknpd.WithToken(expectedToken),
		)

		_, err := c.RequestWithAuth[any](t.Context(), "/error", http.MethodPost, nil)
		if !strings.Contains(err.Error(), "errorCode (502): someMsg") {
			t.Errorf("expected error %q but got %v", "errorCode (502): someMsg", err)
		}
	})

	t.Run("unauthorized error", func(t *testing.T) {
		c := lknpd.New(lknpd.WithBaseURL(testServer.URL))

		_, err := c.RequestWithAuth[any](t.Context(), "/", http.MethodGet, nil)
		if !errors.Is(err, lknpd.ErrUnauthorized) {
			t.Errorf("expected lknpd.ErrUnauthorized but got %v", err)
		}
	})

	t.Run("refresh before request", func(t *testing.T) {
		c := lknpd.New(
			lknpd.WithBaseURL(testServer.URL),
			lknpd.WithRefreshToken(oldRefreshToken),
		)

		if err := c.Refresh(t.Context()); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_, err := c.RequestWithAuth[any](t.Context(), "/", http.MethodGet, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if c.Token() != expectedNewToken {
			t.Errorf("expected c.Token %q but got %q", expectedNewToken, c.Token())
		}

		if c.RefreshToken() != expectedNewRefreshToken {
			t.Errorf("expected c.RefreshToken %q but got %q", expectedNewRefreshToken, c.RefreshToken())
		}

		if !c.TokenExpiresIn().Equal(expectedExpireIn) {
			t.Errorf("expected c.TokenExpireIn %q but got %q", expectedExpireIn.Format(time.DateTime), c.TokenExpiresIn().Format(time.DateTime))
		}
	})

	t.Run("refresh before request when token expired", func(t *testing.T) {
		c := lknpd.New(
			lknpd.WithBaseURL(testServer.URL),
			lknpd.WithRefreshToken(expectedRefreshToken),
		)

		if err := c.Refresh(t.Context()); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_, err := c.RequestWithAuth[any](t.Context(), "/", http.MethodGet, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if c.Token() != expectedNewToken {
			t.Errorf("expected c.Token %q but got %q", expectedNewToken, c.Token())
		}

		if c.RefreshToken() != expectedNewRefreshToken {
			t.Errorf("expected c.RefreshToken %q but got %q", expectedNewRefreshToken, c.RefreshToken())
		}

		if !c.TokenExpiresIn().Equal(expectedExpireIn) {
			t.Errorf("expected c.TokenExpireIn %q but got %q", expectedExpireIn.Format(time.DateTime), c.TokenExpiresIn().Format(time.DateTime))
		}
	})
}

func TestLogin(t *testing.T) {
	expectedINN := "01234"
	expectedPassword := "somePassword"
	expectedNewRefreshToken := "newRefreshToken"
	expectedNewToken := "newToken"
	expectedExpireIn := time.Now().Add(5 * time.Minute).Round(0)

	expectedResponse := lknpd.LoginResponse{
		Token:         expectedNewToken,
		RefreshToken:  expectedNewRefreshToken,
		TokenExpireIn: expectedExpireIn,
		Profile: lknpd.ProfileResponse{
			INN: expectedINN,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/lkfl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResponse)
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	t.Run("successful", func(t *testing.T) {
		c := lknpd.New(
			lknpd.WithBaseURL(testServer.URL),
			lknpd.WithCredentials(expectedINN, expectedPassword),
		)
		resp, err := c.Login(t.Context())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(resp, expectedResponse) {
			t.Errorf("expected response %v but got %v", expectedResponse, resp)
		}

		if c.INN() != expectedINN {
			t.Errorf("expected c.INN %q but got %q", expectedINN, c.INN())
		}
	})
}
