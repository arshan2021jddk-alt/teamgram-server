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
	"github.com/teamgram/proto/mtproto"
	chatpb "github.com/teamgram/teamgram-server/app/service/biz/chat/chat"
	userpb "github.com/teamgram/teamgram-server/app/service/biz/user/user"
)

// ChannelsUpdateUsername
// channels.updateUsername#3514b3de channel:InputChannel username:string = Bool;
func (c *UsernamesCore) ChannelsUpdateUsername(in *mtproto.TLChannelsUpdateUsername) (*mtproto.Bool, error) {
	if in.GetChannel() == nil || in.GetChannel().GetChannelId() == 0 {
		return nil, mtproto.ErrChannelInvalid
	}

	channelId := in.GetChannel().GetChannelId()
	newUsername := in.GetUsername()

	chat, err := c.svcCtx.Dao.ChatClient.ChatGetChatBySelfId(c.ctx, &chatpb.TLChatGetChatBySelfId{
		SelfId: c.MD.UserId,
		ChatId: channelId,
	})
	if err != nil {
		c.Logger.Errorf("channels.updateUsername - get chat error: %v", err)
		return nil, err
	}

	oldUsername := chat.GetUsername()
	if oldUsername == newUsername {
		return mtproto.BoolTrue, nil
	}

	if newUsername != "" {
		ok, err := c.svcCtx.Dao.UserClient.UserUpdateUsernameByUsername(c.ctx, &userpb.TLUserUpdateUsernameByUsername{
			PeerType: mtproto.PEER_CHANNEL,
			PeerId:   channelId,
			Username: newUsername,
		})
		if err != nil {
			return nil, err
		}
		if !mtproto.FromBool(ok) {
			return nil, mtproto.ErrUsernameOccupied
		}
	}

	if oldUsername != "" {
		_, err = c.svcCtx.Dao.UserClient.UserDeleteUsernameByPeer(c.ctx, &userpb.TLUserDeleteUsernameByPeer{
			PeerType: mtproto.PEER_CHANNEL,
			PeerId:   channelId,
		})
		if err != nil {
			c.Logger.Errorf("channels.updateUsername - delete old username error: %v", err)
			return nil, err
		}
	}

	return mtproto.BoolTrue, nil
}
