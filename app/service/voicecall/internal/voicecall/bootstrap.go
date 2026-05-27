package voicecall

import "fmt"

// NewSelfHostedService wires signaling call-state management to a self-hosted LiveKit provider.
// It is intended to be called from the Go signaling service bootstrap.
func NewSelfHostedService(maxParticipants int) (*Service, error) {
	provider, err := NewLiveKitProviderFromEnv()
	if err != nil {
		return nil, fmt.Errorf("init livekit provider: %w", err)
	}
	return NewService(provider, maxParticipants), nil
}
