# Voice Call Service (Group)

This module adds **real group voice call signaling + media integration hooks** for Teamgram using:

- **Go signaling service** (this module)
- **self-hosted LiveKit SFU** (inside Teamgram infra)
- **WebRTC clients** (Android/iOS/Desktop)

## What is implemented

- Group call lifecycle state machine:
  - create call
  - join call
  - leave call
  - discard call
- Participant tracking with role/permission checks.
- Capacity limit enforcement (default 20 users).
- **Real media transport integration point** via SFU provider interface.
- Built-in provider implementation for **LiveKit** token provisioning (JWT based).
- Provider config validation that rejects `*.livekit.cloud` endpoints to enforce self-hosted deployment.

> Important: this service is signaling/control-plane only. Audio media is transported directly between clients and LiveKit SFU.

## Architecture

1. Client captures mic, encodes Opus, publishes/receives tracks with WebRTC.
2. Teamgram signaling calls this service for call state and permissions.
3. Service returns SFU connection metadata (`url`, `room`, `participant`, `token`).
4. Clients connect to self-hosted LiveKit and stream audio in real-time.

## Self-hosted LiveKit deployment

Files:

- `deploy.livekit.yaml` - docker compose stack for self-hosted LiveKit + Redis
- `livekit.yaml` - LiveKit server config (ports, Redis, API keys)

Start local stack:

```bash
docker compose -f app/service/voicecall/deploy.livekit.yaml up -d
```

Set signaling env vars:

- Optional: `LIVEKIT_URL` (default: `ws://127.0.0.1:7880`)
- Optional: `LIVEKIT_API_KEY` (default: `LK_TEAMGRAM_DEV`)
- Optional: `LIVEKIT_API_SECRET` (default: `LK_TEAMGRAM_DEV_SECRET_REPLACE_ME`)

If env vars are not provided, service uses local self-hosted defaults (no cloud account needed).
Use `NewLiveKitProviderFromEnv()` (or `NewSelfHostedService`) in signaling bootstrap.

## Target

- Up to 20 concurrent users per group call (configurable).

## Integration step

Wire these methods into MTProto handlers:

- `phone.createGroupCall`
- `phone.joinGroupCall`
- `phone.leaveGroupCall`
- `phone.discardGroupCall`

and fan out state updates to participants.


## Bootstrap example

```go
svc, err := voicecall.NewSelfHostedService(20)
if err != nil {
    // fail fast: self-hosted LiveKit config is invalid
}
_ = svc
```


## Runnable signaling API (for immediate end-to-end testing)

Run:

```bash
VOICECALL_LISTEN=:18080 go run ./app/service/voicecall/cmd/voicecall
```

Endpoints (`POST` JSON):

- `/voicecall/create` `{ "peer_id": 1001, "user_id": 1 }`
- `/voicecall/join` `{ "peer_id": 1001, "user_id": 2 }` (returns LiveKit token/url/room)
- `/voicecall/leave` `{ "peer_id": 1001, "user_id": 2 }`
- `/voicecall/discard` `{ "peer_id": 1001, "user_id": 1 }`

This allows clients to get real self-hosted LiveKit credentials and start WebRTC audio immediately while MTProto handler wiring is completed.
