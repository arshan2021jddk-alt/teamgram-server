package core

import (
	"testing"

	"github.com/teamgram/proto/mtproto"
)

func TestValidateBotCommandsOK(t *testing.T) {
	err := validateBotCommands([]*mtproto.BotCommand{
		mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "Start", Description: "Start bot"}).To_BotCommand(),
		mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "help_1", Description: "Show help"}).To_BotCommand(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateBotCommandsRejectsInvalid(t *testing.T) {
	cases := []struct {
		name string
		cmds []*mtproto.BotCommand
	}{
		{"nil command", []*mtproto.BotCommand{nil}},
		{"empty name", []*mtproto.BotCommand{mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "", Description: "x"}).To_BotCommand()}},
		{"bad charset", []*mtproto.BotCommand{mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "bad-name", Description: "x"}).To_BotCommand()}},
		{"duplicate", []*mtproto.BotCommand{
			mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "start", Description: "x"}).To_BotCommand(),
			mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "START", Description: "y"}).To_BotCommand(),
		}},
		{"empty description", []*mtproto.BotCommand{mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: "start", Description: ""}).To_BotCommand()}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBotCommands(tt.cmds); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestMakeBotCommandsDOList(t *testing.T) {
	cmds := []*mtproto.BotCommand{
		mtproto.MakeTLBotCommand(&mtproto.BotCommand{Command: " Start ", Description: " Desc "}).To_BotCommand(),
	}
	out := makeBotCommandsDOList(99, cmds)
	if len(out) != 1 {
		t.Fatalf("expected 1 row, got %d", len(out))
	}
	if out[0].BotId != 99 || out[0].Command != "start" || out[0].Description != "Desc" {
		t.Fatalf("unexpected normalized row: %#v", out[0])
	}
}
