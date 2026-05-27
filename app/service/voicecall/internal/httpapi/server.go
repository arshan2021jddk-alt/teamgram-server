package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamgram/teamgram-server/app/service/voicecall/internal/voicecall"
)

type Server struct {
	svc *voicecall.Service
}

func NewServer(svc *voicecall.Service) *Server { return &Server{svc: svc} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/voicecall/create", s.create)
	mux.HandleFunc("/voicecall/join", s.join)
	mux.HandleFunc("/voicecall/leave", s.leave)
	mux.HandleFunc("/voicecall/discard", s.discard)
	return mux
}

type callReq struct {
	PeerID int64 `json:"peer_id"`
	UserID int64 `json:"user_id"`
}

func decode(r *http.Request, out any) error { return json.NewDecoder(r.Body).Decode(out) }
func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { write(w, 405, map[string]string{"error": "method_not_allowed"}); return }
	var req callReq
	if err := decode(r, &req); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	call, err := s.svc.CreateCall(req.PeerID, req.UserID)
	if err != nil { write(w, 500, map[string]string{"error": err.Error()}); return }
	write(w, 200, call)
}

func (s *Server) join(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { write(w, 405, map[string]string{"error": "method_not_allowed"}); return }
	var req callReq
	if err := decode(r, &req); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	call, conn, err := s.svc.JoinCall(req.PeerID, req.UserID)
	if err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	write(w, 200, map[string]any{"call": call, "media": conn})
}

func (s *Server) leave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { write(w, 405, map[string]string{"error": "method_not_allowed"}); return }
	var req callReq
	if err := decode(r, &req); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	if err := s.svc.LeaveCall(req.PeerID, req.UserID); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	write(w, 200, map[string]string{"status": "ok"})
}

func (s *Server) discard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { write(w, 405, map[string]string{"error": "method_not_allowed"}); return }
	var req callReq
	if err := decode(r, &req); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	if err := s.svc.DiscardCall(req.PeerID, req.UserID); err != nil { write(w, 400, map[string]string{"error": err.Error()}); return }
	write(w, 200, map[string]string{"status": "ok"})
}

func ParseMaxParticipants(v string) int {
	n, _ := strconv.Atoi(v)
	if n <= 0 {
		return 20
	}
	return n
}
