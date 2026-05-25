// Copyright (c) 2026 The Teamgram Authors (https://teamgram.net).
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

package core

import (
	"github.com/teamgram/proto/mtproto"
	"github.com/zeromicro/go-zero/core/timex"
)

// MessagesEditChatParticipantRank
// messages.editChatParticipantRank#a00f32b0 peer:InputPeer participant:InputPeer rank:string = Updates;
func (c *ChatsCore) MessagesEditChatParticipantRank(in *mtproto.TLMessagesEditChatParticipantRank) (*mtproto.Updates, error) {
	if in.GetPeer() == nil || in.GetParticipant() == nil {
		err := mtproto.ErrPeerIdInvalid
		c.Logger.Errorf("messages.editChatParticipantRank - invalid peer/participant: %v", err)
		return nil, err
	}

	// Community-compatible fast path: accept request and return empty updates.
	return mtproto.MakeTLUpdates(&mtproto.Updates{
		Updates: []*mtproto.Update{},
		Users:   []*mtproto.User{},
		Chats:   []*mtproto.Chat{},
		Date:    int32(timex.Now()),
		Seq:     0,
	}).To_Updates(), nil
}
