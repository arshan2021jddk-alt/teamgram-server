package voicecall

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrCallNotFound     = errors.New("group call not found")
	ErrCallDiscarded    = errors.New("group call discarded")
	ErrParticipantLimit = errors.New("group call participant limit reached")
)

type Service struct {
	mu             sync.RWMutex
	callsByPeerID  map[int64]*GroupCall
	maxParticipants int
	provider       MediaProvider
}

func NewService(provider MediaProvider, maxParticipants int) *Service {
	if maxParticipants <= 0 {
		maxParticipants = 20
	}
	return &Service{callsByPeerID: map[int64]*GroupCall{}, maxParticipants: maxParticipants, provider: provider}
}

func (s *Service) CreateCall(peerID, creatorUserID int64) (*GroupCall, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	call := &GroupCall{
		CallID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		PeerID:        peerID,
		CreatedByUser: creatorUserID,
		Status:        CallStatusActive,
		Participants:  map[int64]*Participant{},
		CreatedAt:     time.Now(),
	}
	call.Participants[creatorUserID] = &Participant{UserID: creatorUserID, Role: RoleSpeaker, JoinedAt: time.Now(), CanSpeak: true, CanListen: true}
	s.callsByPeerID[peerID] = call
	return call, nil
}

func (s *Service) JoinCall(peerID, userID int64) (*GroupCall, *MediaConnectionInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	call := s.callsByPeerID[peerID]
	if call == nil {
		return nil, nil, ErrCallNotFound
	}
	if call.Status != CallStatusActive {
		return nil, nil, ErrCallDiscarded
	}
	if _, ok := call.Participants[userID]; !ok {
		if len(call.Participants) >= s.maxParticipants {
			return nil, nil, ErrParticipantLimit
		}
		call.Participants[userID] = &Participant{UserID: userID, Role: RoleSpeaker, JoinedAt: time.Now(), CanSpeak: true, CanListen: true}
	}
	conn, err := s.provider.IssueJoinCredential(call, userID)
	if err != nil {
		return nil, nil, err
	}
	return call, conn, nil
}

func (s *Service) LeaveCall(peerID, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	call := s.callsByPeerID[peerID]
	if call == nil {
		return ErrCallNotFound
	}
	delete(call.Participants, userID)
	if len(call.Participants) == 0 {
		call.Status = CallStatusDiscarded
		delete(s.callsByPeerID, peerID)
	}
	return nil
}

func (s *Service) DiscardCall(peerID, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	call := s.callsByPeerID[peerID]
	if call == nil {
		return ErrCallNotFound
	}
	if call.CreatedByUser != userID {
		return errors.New("only creator can discard group call")
	}
	call.Status = CallStatusDiscarded
	delete(s.callsByPeerID, peerID)
	return nil
}
