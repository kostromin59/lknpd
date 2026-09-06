package lknpd

const (
	defaultAppVersion = "1.0.0"
	defaultSourceType = "WEB"
	defaultBaseURL    = "https://lknpd.nalog.ru"
	defaultUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36"
)

type options struct {
	appVersion string
	sourceType string
	deviceID   string
	userAgent  string
	baseURL    string
}

type option func(options) options

func newOptions() options {
	return options{
		appVersion: defaultAppVersion,
		sourceType: defaultAppVersion,
		deviceID:   generateDeviceID(),
		userAgent:  defaultUserAgent,
		baseURL:    defaultBaseURL,
	}
}

func WithAppVersion(appVersion string) option {
	return func(o options) options {
		o.appVersion = appVersion
		return o
	}
}

func WithSourceType(sourceType string) option {
	return func(o options) options {
		o.sourceType = sourceType
		return o
	}
}

func WithDeviceID(deviceID string) option {
	return func(o options) options {
		o.deviceID = deviceID
		return o
	}
}

func WithUserAgent(userAgent string) option {
	return func(o options) options {
		o.userAgent = userAgent
		return o
	}
}

func WithBaseURL(baseURL string) option {
	return func(o options) options {
		o.baseURL = baseURL
		return o
	}
}
