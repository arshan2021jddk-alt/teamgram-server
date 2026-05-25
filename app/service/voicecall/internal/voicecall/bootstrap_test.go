package voicecall

import (
	"os"
	"testing"
)

func TestNewSelfHostedService(t *testing.T) {
	t.Setenv(EnvLiveKitURL, "wss://voice.internal")
	t.Setenv(EnvLiveKitAPIKey, "k")
	t.Setenv(EnvLiveKitAPISecret, "s")

	svc, err := NewSelfHostedService(20)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
	}
}

func TestNewSelfHostedServiceFailsWithoutEnv(t *testing.T) {
	_ = os.Unsetenv(EnvLiveKitURL)
	_ = os.Unsetenv(EnvLiveKitAPIKey)
	_ = os.Unsetenv(EnvLiveKitAPISecret)
	if _, err := NewSelfHostedService(20); err != nil {
		t.Fatalf("unexpected error with defaults: %v", err)
	}
}
