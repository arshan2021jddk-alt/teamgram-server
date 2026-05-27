package voicecall

import "testing"

type fakeProvider struct{}

func (f *fakeProvider) IssueJoinCredential(call *GroupCall, userID int64) (*MediaConnectionInfo, error) {
	return &MediaConnectionInfo{Provider: "fake", ServerURL: "wss://voice.internal", RoomName: call.CallID, Participant: "u"}, nil
}

func TestJoinCallParticipantLimit(t *testing.T) {
	svc := NewService(&fakeProvider{}, 2)
	_, err := svc.CreateCall(10, 100)
	if err != nil {
		t.Fatalf("create call: %v", err)
	}
	_, _, err = svc.JoinCall(10, 101)
	if err != nil {
		t.Fatalf("join #2 unexpected err: %v", err)
	}
	_, _, err = svc.JoinCall(10, 102)
	if err != ErrParticipantLimit {
		t.Fatalf("want ErrParticipantLimit, got: %v", err)
	}
}

func TestDiscardOnlyCreator(t *testing.T) {
	svc := NewService(&fakeProvider{}, 20)
	_, _ = svc.CreateCall(10, 100)
	if err := svc.DiscardCall(10, 101); err == nil {
		t.Fatal("expected error when non-creator discards")
	}
}
