package voicecall

import "testing"

func TestLiveKitValidateRejectsCloud(t *testing.T) {
	p := &LiveKitProvider{ServerURL: "wss://abc.livekit.cloud", APIKey: "k", APISecret: "s"}
	err := p.Validate()
	if err != ErrLiveKitCloudDisallowed {
		t.Fatalf("want ErrLiveKitCloudDisallowed, got: %v", err)
	}
}

func TestLiveKitValidateSelfHosted(t *testing.T) {
	p := &LiveKitProvider{ServerURL: "wss://livekit.voice.svc.cluster.local", APIKey: "k", APISecret: "s"}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected validate err: %v", err)
	}
}
