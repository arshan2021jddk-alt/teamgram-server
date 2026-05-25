package voicecall

import "time"

type CallStatus string

const (
	CallStatusActive    CallStatus = "active"
	CallStatusDiscarded CallStatus = "discarded"
)

type ParticipantRole string

const (
	RoleSpeaker  ParticipantRole = "speaker"
	RoleListener ParticipantRole = "listener"
)

type GroupCall struct {
	CallID        string
	PeerID        int64
	CreatedByUser int64
	Status        CallStatus
	Participants  map[int64]*Participant
	CreatedAt     time.Time
}

type Participant struct {
	UserID    int64
	Role      ParticipantRole
	JoinedAt  time.Time
	Muted     bool
	CanSpeak  bool
	CanListen bool
}

type MediaConnectionInfo struct {
	Provider    string
	ServerURL   string
	RoomName    string
	Participant string
	Token       string
}
