package lknpd

const (
	defaultAppVersion = "1.0.0"
	defaultSourceType = "WEB"
)

type DeviceInfoRequest struct {
	AppVersion     string             `json:"appVersion,omitempty"`
	SourceDeviceID string             `json:"sourceDeviceId,omitempty"`
	SourceType     string             `json:"sourceType,omitempty"`
	MetaDetails    MetaDetailsRequset `json:"metaDetails,omitempty"`
}

type MetaDetailsRequset struct {
	UserAgent string `json:"userAgent,omitempty"`
}

func NewDeviceInfoRequest(deviceID string, userAgent string) DeviceInfoRequest {
	return DeviceInfoRequest{
		AppVersion:     defaultAppVersion,
		SourceType:     defaultSourceType,
		SourceDeviceID: deviceID,
		MetaDetails: MetaDetailsRequset{
			UserAgent: userAgent,
		},
	}
}
