/*
 * Created from 'scheme.tl' by 'mtprotoc'
 *
 * Copyright (c) 2021-present,  Teamgram Studio (https://teamgram.io).
 *  All rights reserved.
 *
 * Author: teamgramio (teamgram.io@gmail.com)
 */

package core

import (
	"regexp"
	"strings"

	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/biz/user/internal/dal/dataobject"
	"github.com/teamgram/teamgram-server/app/service/biz/user/user"
)

const (
	maxBotCommands           = 100
	maxBotCommandNameLen     = 32
	maxBotCommandDescLenRune = 256
)

var botCommandPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

func normalizeBotCommandName(v string) string {
	return strings.TrimSpace(strings.ToLower(v))
}

func normalizeBotCommandDesc(v string) string {
	return strings.TrimSpace(v)
}

func validateBotCommands(commands []*mtproto.BotCommand) error {
	if len(commands) > maxBotCommands {
		return mtproto.ErrBadRequest
	}

	seen := make(map[string]struct{}, len(commands))
	for _, cmd := range commands {
		if cmd == nil {
			return mtproto.ErrBadRequest
		}

		name := normalizeBotCommandName(cmd.GetCommand())
		desc := normalizeBotCommandDesc(cmd.GetDescription())
		if name == "" || len(name) > maxBotCommandNameLen || !botCommandPattern.MatchString(name) {
			return mtproto.ErrBadRequest
		}
		if desc == "" || len([]rune(desc)) > maxBotCommandDescLenRune {
			return mtproto.ErrBadRequest
		}
		if _, ok := seen[name]; ok {
			return mtproto.ErrBadRequest
		}
		seen[name] = struct{}{}
	}

	return nil
}

func makeBotCommandsDOList(botID int64, commands []*mtproto.BotCommand) []*dataobject.BotCommandsDO {
	bulk := make([]*dataobject.BotCommandsDO, 0, len(commands))
	for _, cmd := range commands {
		bulk = append(bulk, &dataobject.BotCommandsDO{
			BotId:       botID,
			Command:     normalizeBotCommandName(cmd.GetCommand()),
			Description: normalizeBotCommandDesc(cmd.GetDescription()),
		})
	}

	return bulk
}

// UserSetBotCommands
// user.setBotCommands user_id:long bot_id:long commands:Vector<BotCommand> = Bool;
func (c *UserCore) UserSetBotCommands(in *user.TLUserSetBotCommands) (*mtproto.Bool, error) {
	if in == nil || in.BotId <= 0 {
		return nil, mtproto.ErrBadRequest
	}
	if c.MD == nil || c.MD.UserId == 0 {
		return nil, mtproto.ErrBotInvalid
	}
	if in.UserId > 0 && c.MD.UserId != in.UserId {
		return nil, mtproto.ErrBotInvalid
	}

	if err := validateBotCommands(in.Commands); err != nil {
		return nil, err
	}

	botsDO, err := c.svcCtx.BotsDAO.Select(c.ctx, in.BotId)
	if err != nil {
		return nil, err
	}
	if botsDO == nil || botsDO.CreatorUserId != c.MD.UserId {
		return nil, mtproto.ErrBotInvalid
	}

	tR := sqlx.TxWrapper(c.ctx, c.svcCtx.Dao.DB, func(tx *sqlx.Tx, result *sqlx.StoreResult) {
		_, result.Err = c.svcCtx.Dao.BotCommandsDAO.DeleteTx(tx, in.BotId)
		if result.Err != nil {
			return
		}

		if len(in.Commands) == 0 {
			return
		}

		bulk := makeBotCommandsDOList(in.BotId, in.Commands)
		_, _, result.Err = c.svcCtx.Dao.BotCommandsDAO.InsertBulkTx(tx, bulk)
	})
	if tR.Err != nil {
		return nil, tR.Err
	}

	return mtproto.BoolTrue, nil
}
