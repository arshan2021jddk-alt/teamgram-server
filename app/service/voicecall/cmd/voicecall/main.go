package main

import (
	"log"
	"net/http"
	"os"

	httpapi "github.com/teamgram/teamgram-server/app/service/voicecall/internal/httpapi"
	"github.com/teamgram/teamgram-server/app/service/voicecall/internal/voicecall"
)

func main() {
	svc, err := voicecall.NewSelfHostedService(httpapi.ParseMaxParticipants(os.Getenv("VOICECALL_MAX_PARTICIPANTS")))
	if err != nil {
		log.Fatalf("voicecall init failed: %v", err)
	}
	addr := os.Getenv("VOICECALL_LISTEN")
	if addr == "" {
		addr = ":18080"
	}
	log.Printf("voicecall signaling listening on %s", addr)
	if err = http.ListenAndServe(addr, httpapi.NewServer(svc).Handler()); err != nil {
		log.Fatal(err)
	}
}
