package lknpd

const (
	defaultAppVersion = "1.0.0"
	defaultSourceType = "WEB"
	defaultBaseURL    = "https://lknpd.nalog.ru"
	defaultUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36"
)

type options struct {
	inn          string
	password     string
	appVersion   string
	sourceType   string
	deviceID     string
	userAgent    string
	baseURL      string
	token        string
	refreshToken string
}

type Option func(options) options

func newOptions() options {
	return options{
		appVersion: defaultAppVersion,
		sourceType: defaultAppVersion,
		deviceID:   GenerateDeviceID(),
		userAgent:  defaultUserAgent,
		baseURL:    defaultBaseURL,
	}
}

// Default value: "1.0.0"
func WithAppVersion(appVersion string) Option {
	return func(o options) options {
		o.appVersion = appVersion
		return o
	}
}

// Default value: "WEB"
func WithSourceType(sourceType string) Option {
	return func(o options) options {
		o.sourceType = sourceType
		return o
	}
}

// Default value: generated string using [GenerateDeviceID].
func WithDeviceID(deviceID string) Option {
	return func(o options) options {
		o.deviceID = deviceID
		return o
	}
}

// Default value: generated string using [GenerateDeviceID].
func WithUserAgent(userAgent string) Option {
	return func(o options) options {
		o.userAgent = userAgent
		return o
	}
}

// Default value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36"
func WithBaseURL(baseURL string) Option {
	return func(o options) options {
		o.baseURL = baseURL
		return o
	}
}

// Default value: empty. Use login to get tokens (inn and password are required).
func WithToken(token string) Option {
	return func(o options) options {
		o.token = token
		return o
	}
}

// Default value: empty. Use login to get tokens (inn and password are required).
func WithRefreshToken(refreshToken string) Option {
	return func(o options) options {
		o.refreshToken = refreshToken
		return o
	}
}

// Default value: empty. Required to get tokens. May be combined with tokens ([WithTokens])
func WithCredentials(inn, password string) Option {
	return func(o options) options {
		o.inn = inn
		o.password = password
		return o
	}
}
