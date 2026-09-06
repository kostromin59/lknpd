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

	token          string
	tokenExpiresIn time.Time
	refreshToken   string

	mu *sync.RWMutex
}

// New creates new instance of [Client]. When the provided baseURL is empty, `https://lknpd.nalog.ru` will be used.
func New(opts ...option) *Client {
	o := newOptions()

	for _, opt := range opts {
		o = opt(o)
	}

	return &Client{
		c:       new(http.Client),
		baseURL: o.baseURL,
		deviceInfo: DeviceInfo{
			AppVersion:     o.appVersion,
			SourceDeviceID: o.deviceID,
			SourceType:     o.sourceType,
			MetaDetails: MetaDetails{
				UserAgent: o.userAgent,
			},
		},
		mu: new(sync.RWMutex),
	}
}

func (c *Client) Login(ctx context.Context, inn, password string) (LoginResponse, error) {
	const op = "lknpd.Client.Login"

	c.mu.RLock()
	body := LoginRequest{
		Username:   inn,
		Password:   password,
		DeviceInfo: c.deviceInfo,
	}
	c.mu.RUnlock()

	resp, err := c.Request[LoginResponse](ctx, "/api/v1/auth/lkfl", http.MethodPost, body)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	c.mu.Lock()
	c.token = resp.Token
	c.tokenExpiresIn = resp.TokenExpiresIn
	c.refreshToken = resp.RefreshToken
	c.mu.Unlock()

	return resp, nil
}

func (c *Client) CreateIncome(ctx context.Context, client IncomeClient, services []Income, date time.Time) (string, error) {
	const op = "lknpd.Client.CreateIncome"

	total := .0
	for _, income := range services {
		total += income.Amount * float64(income.Quantity)
	}

	body := CreateIncomeRequest{
		Client:                           client,
		IgnoreMaxTotalIncomeRestrictions: false,
		OperationTime:                    date,
		PaymentType:                      defaultPaymentType,
		RequestTime:                      time.Now(),
		Services:                         services,
		TotalAmount:                      fmt.Sprintf("%.2f", total),
	}

	resp, err := c.Request[CreateIncomeResponse](ctx, "/api/v1/income", http.MethodPost, body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ApprovedReceiptUUID, nil
}

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

	resp, err := c.c.Do(req)
	if err != nil {
		return responseValue, fmt.Errorf("%s: %w", op, err)
	}

	// Bad Status Code
	if resp.StatusCode > 299 {
		defer resp.Body.Close()

		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return responseValue, fmt.Errorf("%s: %w", op, err)
		}

		return responseValue, fmt.Errorf("%s: %s (%d): %s", op, errorResponse.Code, resp.StatusCode, errorResponse.Message)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&responseValue); err != nil {
			return responseValue, fmt.Errorf("%s: %w", op, err)
		}
	}

	return responseValue, nil
}

func (c *Client) RequestWithAuth[T any](ctx context.Context, path, method string, body any) (T, error) {
	var zero T

	if c.tokenExpiresIn.After(time.Now()) {
		if err := c.RefreshToken(ctx); err != nil {
			return zero, err
		}
	}

	return c.Request[T](ctx, path, method, body)
}

func (c *Client) RefreshToken(ctx context.Context) error {
	const op = "lknpd.Client.RefreshToken"

	c.mu.RLock()
	body := RefreshTokenRequest{
		DeviceInfo:   c.deviceInfo,
		RefreshToken: c.refreshToken,
	}
	c.mu.RUnlock()

	resp, err := c.Request[RefreshTokenResponse](ctx, "/api/v1/auth/token", http.MethodPost, body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.mu.Lock()
	c.token = resp.Token
	c.tokenExpiresIn = resp.TokenExpiresIn
	if resp.RefreshToken != "" {
		c.refreshToken = resp.RefreshToken
	}
	c.mu.Unlock()

	return nil
}
