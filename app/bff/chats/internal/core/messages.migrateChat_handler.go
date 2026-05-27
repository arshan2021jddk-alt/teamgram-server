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
	"math/rand"

	"github.com/teamgram/proto/mtproto"
	chatpb "github.com/teamgram/teamgram-server/app/service/biz/chat/chat"
	"github.com/zeromicro/go-zero/core/timex"
)

// MessagesMigrateChat
// messages.migrateChat#a2875319 chat_id:long = Updates;
func (c *ChatsCore) MessagesMigrateChat(in *mtproto.TLMessagesMigrateChat) (*mtproto.Updates, error) {
	if in.GetChatId() <= 0 {
		err := mtproto.ErrChatIdInvalid
		c.Logger.Errorf("messages.migrateChat - invalid chat_id: %v", err)
		return nil, err
	}

	mChat, err := c.svcCtx.Dao.ChatClient.Client().ChatGetMutableChat(c.ctx, &chatpb.TLChatGetMutableChat{
		ChatId: in.GetChatId(),
	})
	if err != nil {
		c.Logger.Errorf("messages.migrateChat - get mutable chat error: %v", err)
		return nil, err
	}

	channelId := c.svcCtx.Dao.NextId(c.ctx)
	if channelId <= 0 {
		err = mtproto.ErrInternal
		c.Logger.Errorf("messages.migrateChat - next id error: %v", err)
		return nil, err
	}
	accessHash := rand.Int63()

	_, err = c.svcCtx.Dao.ChatClient.ChatMigratedToChannel(c.ctx, &chatpb.TLChatMigratedToChannel{
		Chat:       mChat,
		Id:         channelId,
		AccessHash: accessHash,
	})
	if err != nil {
		c.Logger.Errorf("messages.migrateChat - migrate error: %v", err)
		return nil, err
	}

	// Return empty updates for now; migration state is persisted in chat service.
	return mtproto.MakeTLUpdates(&mtproto.Updates{
		Updates: []*mtproto.Update{},
		Users:   []*mtproto.User{},
		Chats:   []*mtproto.Chat{},
		Date:    int32(timex.Now()),
		Seq:     0,
	}).To_Updates(), nil
}
