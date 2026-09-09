package lknpd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	c          *http.Client
	baseURL    string
	deviceInfo DeviceInfo

	inn      string
	password string

	token         string
	refreshToken  string
	tokenExpireIn *time.Time

	mu *sync.RWMutex
}

// New creates new instance of [Client]. Use options to change baseURL and another settings.
func New(opts ...Option) *Client {
	o := newOptions()

	for _, opt := range opts {
		o = opt(o)
	}

	deviceID := o.deviceID
	if o.tokenDeviceID != "" {
		deviceID = o.tokenDeviceID
	}

	if o.refreshDeviceID != "" {
		deviceID = o.refreshDeviceID
	}

	return &Client{
		inn:      o.inn,
		password: o.password,
		c:        new(http.Client),
		baseURL:  o.baseURL,
		deviceInfo: DeviceInfo{
			AppVersion:     o.appVersion,
			SourceDeviceID: deviceID,
			SourceType:     o.sourceType,
			MetaDetails: MetaDetails{
				UserAgent: o.userAgent,
			},
		},
		token:        o.token,
		refreshToken: o.refreshToken,
		mu:           new(sync.RWMutex),
	}
}

// INN returns inn provided in [New].
func (c *Client) INN() string {
	c.mu.RLock()
	inn := c.inn
	c.mu.RUnlock()

	return inn
}

// DeviceInfo returns [DeviceInfo]. It stores generated deviceID.
func (c *Client) DeviceInfo() DeviceInfo {
	c.mu.RLock()
	deviceInfo := c.deviceInfo
	c.mu.RUnlock()

	return deviceInfo
}

func (c *Client) Token() string {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	return token
}

func (c *Client) TokenExpiresIn() *time.Time {
	c.mu.RLock()
	tokenExpiresIn := *c.tokenExpireIn
	c.mu.RUnlock()

	return &tokenExpiresIn
}

func (c *Client) RefreshToken() string {
	c.mu.RLock()
	refreshToken := c.refreshToken
	c.mu.RUnlock()

	return refreshToken
}

// Login requests tokens using inn, password and [DeviceInfo].
func (c *Client) Login(ctx context.Context) (LoginResponse, error) {
	const op = "lknpd.Client.Login"

	c.mu.RLock()
	body := LoginRequest{
		Username:   c.inn,
		Password:   c.password,
		DeviceInfo: c.deviceInfo,
	}
	c.mu.RUnlock()

	if body.Username == "" || body.Password == "" {
		return LoginResponse{}, fmt.Errorf("%s: %w", op, ErrUnauthorized)
	}

	resp, err := c.Request[LoginResponse](ctx, "/api/v1/auth/lkfl", http.MethodPost, body)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	c.mu.Lock()
	c.token = resp.Token
	c.tokenExpireIn = &resp.TokenExpireIn
	c.refreshToken = resp.RefreshToken
	c.inn = resp.Profile.INN
	c.mu.Unlock()

	return resp, nil
}

func (c *Client) Logout(ctx context.Context) error {
	const op = "lknpd.Client.Logout"

	_, err := c.Request[any](ctx, "/api/v1/auth/esia/logout", http.MethodPost, map[string]string{"redirectUrl": "https://lknpd.nalog.ru/auth/login"})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.mu.Lock()
	c.refreshToken = ""
	c.token = ""
	c.inn = ""
	c.tokenExpireIn = nil
	c.mu.Unlock()

	return nil
}

// CreateIncome creates income and returns ApprovedReceiptUUID.
func (c *Client) CreateIncome(ctx context.Context, client IncomeClient, services []Income, date time.Time) (string, error) {
	const op = "lknpd.Client.CreateIncome"

	total := .0
	for _, income := range services {
		total += income.Amount * float64(income.Quantity)
	}

	body := CreateIncomeRequest{
		Client:                           client,
		IgnoreMaxTotalIncomeRestrictions: defaultIgnoreMaxTotalIncomeRestrictions,
		OperationTime:                    date,
		PaymentType:                      defaultPaymentType,
		RequestTime:                      date,
		Services:                         services,
		TotalAmount:                      fmt.Sprintf("%.2f", total),
	}

	resp, err := c.RequestWithAuth[CreateIncomeResponse](ctx, "/api/v1/income", http.MethodPost, body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ApprovedReceiptUUID, nil
}

func (c *Client) CancelIncome(ctx context.Context, receiptUUID string, comment CancelIncomeComment) error {
	const op = "lknpd.Client.CancelIncome"

	now := time.Now()
	body := CancelIncomeRequest{
		Comment:       comment,
		OperationTime: now,
		PartnerCode:   nil,
		ReceiptUUID:   receiptUUID,
		RequestTime:   now,
	}

	_, err := c.RequestWithAuth[any](ctx, "/api/v1/cancel", http.MethodPost, body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Request executes HTTP request and returns response body (use generic type).
func (c *Client) Request[T any](ctx context.Context, path, method string, body any) (T, error) {
	const op = "lknpd.Client.Request"

	c.mu.RLock()
	defer c.mu.RUnlock()

	var responseValue T

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return responseValue, fmt.Errorf("%s: %w", op, err)
	}
	u = u.JoinPath(path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return responseValue, fmt.Errorf("%s: %w", op, err)
		}

		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return responseValue, fmt.Errorf("%s: %w", op, err)
	}

	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.c.Do(req)
	if err != nil {
		return responseValue, fmt.Errorf("%s: %w", op, err)
	}

	// Bad Status Code
	if resp.StatusCode > 299 {
		defer func() {
			_ = resp.Body.Close()
		}()

		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return responseValue, fmt.Errorf("%s: %w", op, err)
		}

		return responseValue, fmt.Errorf("%s: %s (%d): %s", op, errorResponse.Code, resp.StatusCode, errorResponse.Message)
	}

	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()

		if err := json.NewDecoder(resp.Body).Decode(&responseValue); err != nil {
			return responseValue, fmt.Errorf("%s: %w", op, err)
		}
	}

	return responseValue, nil
}

// RequestWithAuth requires at least token or refresh token. If token is empty or is expired it will be refreshed. If both tokens are empty returns [ErrUnauthorized]
func (c *Client) RequestWithAuth[T any](ctx context.Context, path, method string, body any) (T, error) {
	const op = "lknpd.RequestWithAuth"
	var zero T

	c.mu.RLock()
	token := c.token
	refreshToken := c.refreshToken

	var tokenExpireIn *time.Time
	if c.tokenExpireIn != nil {
		tokenExpireIn = new(*c.tokenExpireIn)
	}
	c.mu.RUnlock()

	if token == "" && refreshToken == "" {
		return zero, fmt.Errorf("%s: %w", op, ErrUnauthorized)
	}

	if token == "" {
		if err := c.Refresh(ctx); err != nil {
			return zero, fmt.Errorf("%s: %w", op, err)
		}
	}

	if tokenExpireIn != nil && tokenExpireIn.Before(time.Now()) {
		if err := c.Refresh(ctx); err != nil {
			return zero, fmt.Errorf("%s: %w", op, err)
		}
	}

	resp, err := c.Request[T](ctx, path, method, body)
	if err != nil {
		return zero, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

// Refresh executes HTTP request to refresh tokens.
func (c *Client) Refresh(ctx context.Context) error {
	const op = "lknpd.Client.Refresh"

	c.mu.RLock()
	body := RefreshTokenRequest{
		DeviceInfo:   c.deviceInfo,
		RefreshToken: c.refreshToken,
	}
	c.mu.RUnlock()

	if body.RefreshToken == "" {
		return fmt.Errorf("%s: %w", op, ErrUnauthorized)
	}

	resp, err := c.Request[RefreshTokenResponse](ctx, "/api/v1/auth/token", http.MethodPost, body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.mu.Lock()
	c.token = resp.Token
	c.tokenExpireIn = &resp.TokenExpireIn
	if resp.RefreshToken != "" {
		c.refreshToken = resp.RefreshToken
	}
	c.mu.Unlock()

	return nil
}
