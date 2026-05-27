package core

import (
	"strings"
	"testing"

	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/biz/user/user"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestMakeBotDataUpdateMap(t *testing.T) {
	in := &user.TLUserUpdateBotData{
		BotChatHistory:       mtproto.BoolTrue,
		BotNochats:           mtproto.BoolFalse,
		BotInlineGeo:         mtproto.BoolTrue,
		BotAttachMenu:        mtproto.BoolTrue,
		BotInlinePlaceholder: wrapperspb.String("  hello bot  "),
		BotHasMainApp:        mtproto.BoolFalse,
	}

	m, err := makeBotDataUpdateMap(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 7 {
		t.Fatalf("unexpected map size: %d", len(m))
	}
	if m["bot_chat_history"] != true || m["bot_nochats"] != false || m["bot_inline_geo"] != true {
		t.Fatalf("unexpected bool values: %#v", m)
	}
	if m["bot_attach_menu"] != true || m["attach_menu_enabled"] != true {
		t.Fatalf("attach menu flags not synced: %#v", m)
	}
	if m["bot_inline_placeholder"] != "hello bot" {
		t.Fatalf("placeholder not trimmed: %#v", m["bot_inline_placeholder"])
	}
	if m["bot_has_main_app"] != false {
		t.Fatalf("unexpected bot_has_main_app: %#v", m["bot_has_main_app"])
	}
}

func TestMakeBotDataUpdateMapRejectsLongPlaceholder(t *testing.T) {
	in := &user.TLUserUpdateBotData{BotInlinePlaceholder: wrapperspb.String(strings.Repeat("ا", 65))}
	if _, err := makeBotDataUpdateMap(in); err == nil {
		t.Fatalf("expected error for long placeholder")
	}
}
