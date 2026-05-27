package voicecall

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MediaProvider interface {
	IssueJoinCredential(call *GroupCall, userID int64) (*MediaConnectionInfo, error)
}

type LiveKitProvider struct {
	ServerURL string
	APIKey    string
	APISecret string
}

func (p *LiveKitProvider) IssueJoinCredential(call *GroupCall, userID int64) (*MediaConnectionInfo, error) {
	room := fmt.Sprintf("group-%d-call-%s", call.PeerID, call.CallID)
	identity := fmt.Sprintf("u%d", userID)

	claims := jwt.MapClaims{
		"iss": p.APIKey,
		"sub": identity,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"video": map[string]any{
			"roomJoin": true,
			"room":     room,
			"canPublish": true,
			"canSubscribe": true,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(p.APISecret))
	if err != nil {
		return nil, err
	}

	return &MediaConnectionInfo{
		Provider:    "livekit",
		ServerURL:   p.ServerURL,
		RoomName:    room,
		Participant: identity,
		Token:       signed,
	}, nil
}
