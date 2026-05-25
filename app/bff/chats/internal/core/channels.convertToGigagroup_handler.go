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
	"github.com/zeromicro/go-zero/core/timex"
)

// ChannelsConvertToGigagroup
// channels.convertToGigagroup#b290c69 channel:InputChannel = Updates;
func (c *ChatsCore) ChannelsConvertToGigagroup(in *mtproto.TLChannelsConvertToGigagroup) (*mtproto.Updates, error) {
	if in.GetChannel() == nil || in.GetChannel().GetChannelId() == 0 {
		err := mtproto.ErrChannelInvalid
		c.Logger.Errorf("channels.convertToGigagroup - invalid channel: %v", err)
		return nil, err
	}

	// Community fast-path: accept request and return empty updates until dedicated backend is wired.
	return mtproto.MakeTLUpdates(&mtproto.Updates{
		Updates: []*mtproto.Update{},
		Users:   []*mtproto.User{},
		Chats:   []*mtproto.Chat{},
		Date:    int32(timex.Now()),
		Seq:     0,
	}).To_Updates(), nil
}
