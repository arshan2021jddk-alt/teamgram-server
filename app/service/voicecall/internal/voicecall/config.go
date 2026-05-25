package voicecall

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

const (
	EnvLiveKitURL       = "LIVEKIT_URL"
	EnvLiveKitAPIKey    = "LIVEKIT_API_KEY"
	EnvLiveKitAPISecret = "LIVEKIT_API_SECRET"
	DefaultLiveKitURL   = "ws://127.0.0.1:7880"
	DefaultLiveKitAPIKey = "LK_TEAMGRAM_DEV"
	DefaultLiveKitAPISecret = "LK_TEAMGRAM_DEV_SECRET_REPLACE_ME"
)

var ErrLiveKitCloudDisallowed = errors.New("livekit cloud endpoints are disallowed; use self-hosted LiveKit URL")

func NewLiveKitProviderFromEnv() (*LiveKitProvider, error) {
	provider := &LiveKitProvider{
		ServerURL: strings.TrimSpace(os.Getenv(EnvLiveKitURL)),
		APIKey:    strings.TrimSpace(os.Getenv(EnvLiveKitAPIKey)),
		APISecret: strings.TrimSpace(os.Getenv(EnvLiveKitAPISecret)),
	}
	if provider.ServerURL == "" {
		provider.ServerURL = DefaultLiveKitURL
	}
	if provider.APIKey == "" {
		provider.APIKey = DefaultLiveKitAPIKey
	}
	if provider.APISecret == "" {
		provider.APISecret = DefaultLiveKitAPISecret
	}
	if err := provider.Validate(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *LiveKitProvider) Validate() error {
	if p.ServerURL == "" || p.APIKey == "" || p.APISecret == "" {
		return errors.New("livekit config is incomplete")
	}
	u, err := url.Parse(p.ServerURL)
	if err != nil || u.Host == "" {
		return errors.New("invalid livekit url")
	}
	host := strings.ToLower(u.Hostname())
	if strings.HasSuffix(host, ".livekit.cloud") || host == "livekit.cloud" {
		return ErrLiveKitCloudDisallowed
	}
	return nil
}
