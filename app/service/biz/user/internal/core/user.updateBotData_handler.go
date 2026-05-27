// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package core

import (
	"strings"

	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/biz/user/user"
)

func canMutateBotByOwner(mdUserID int64, creatorUserID int64) bool {
	return mdUserID > 0 && creatorUserID > 0 && mdUserID == creatorUserID
}

func makeBotDataUpdateMap(in *user.TLUserUpdateBotData) (map[string]interface{}, error) {
	cMap := map[string]interface{}{}

	if in.GetBotChatHistory() != nil {
		cMap["bot_chat_history"] = in.GetBotChatHistory().PredicateName == mtproto.Predicate_boolTrue
	}
	if in.GetBotNochats() != nil {
		cMap["bot_nochats"] = in.GetBotNochats().PredicateName == mtproto.Predicate_boolTrue
	}
	if in.GetBotInlineGeo() != nil {
		cMap["bot_inline_geo"] = in.GetBotInlineGeo().PredicateName == mtproto.Predicate_boolTrue
	}
	if in.GetBotAttachMenu() != nil {
		v := in.GetBotAttachMenu().PredicateName == mtproto.Predicate_boolTrue
		cMap["bot_attach_menu"] = v
		cMap["attach_menu_enabled"] = v
	}
	if in.GetBotInlinePlaceholder() != nil {
		placeholder := strings.TrimSpace(in.GetBotInlinePlaceholder().GetValue())
		if len([]rune(placeholder)) > 64 {
			return nil, mtproto.NewRpcError(mtproto.ErrBadRequest)
		}
		cMap["bot_inline_placeholder"] = placeholder
	}
	if in.GetBotHasMainApp() != nil {
		cMap["bot_has_main_app"] = in.GetBotHasMainApp().PredicateName == mtproto.Predicate_boolTrue
	}

	return cMap, nil
}

// UserUpdateBotData
// user.updateBotData flags:# user_id:long bot_chat_history:flags.15?Bool bot_nochats:flags.16?Bool bot_inline_geo:flags.21?Bool bot_attach_menu:flags.27?Bool bot_inline_placeholder:flags.19?string = Bool;
func (c *UserCore) UserUpdateBotData(in *user.TLUserUpdateBotData) (*mtproto.Bool, error) {
	if in == nil || in.BotId <= 0 {
		return nil, mtproto.NewRpcError(mtproto.ErrBadRequest)
	}

	botsDO, err := c.svcCtx.BotsDAO.Select(c.ctx, in.BotId)
	if err != nil {
		return nil, err
	}
	if botsDO == nil {
		return nil, mtproto.ErrBotInvalid
	}
	if c.MD == nil || !canMutateBotByOwner(c.MD.UserId, botsDO.CreatorUserId) {
		return nil, mtproto.ErrBotInvalid
	}

	cMap, err := makeBotDataUpdateMap(in)
	if err != nil {
		return nil, err
	}

	if len(cMap) == 0 {
		return mtproto.BoolTrue, nil
	}

	_, err = c.svcCtx.BotsDAO.Update(c.ctx, cMap, in.BotId)
	if err != nil {
		return nil, err
	}

	return mtproto.BoolTrue, nil
}
