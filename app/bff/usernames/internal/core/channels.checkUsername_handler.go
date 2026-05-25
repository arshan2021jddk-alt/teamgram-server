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
	"github.com/teamgram/marmota/pkg/strings2"
	"github.com/teamgram/marmota/pkg/utils"
	"github.com/teamgram/proto/mtproto"
	userpb "github.com/teamgram/teamgram-server/app/service/biz/user/user"
)

// ChannelsCheckUsername
// channels.checkUsername#10e6bd2c channel:InputChannel username:string = Bool;
func (c *UsernamesCore) ChannelsCheckUsername(in *mtproto.TLChannelsCheckUsername) (*mtproto.Bool, error) {
	if in.GetChannel() == nil || in.GetChannel().GetChannelId() == 0 {
		return nil, mtproto.ErrChannelInvalid
	}

	username := in.GetUsername()
	if len(username) < userpb.MinUsernameLen ||
		!strings2.IsAlNumString(username) ||
		utils.IsNumber(username[0]) {
		return nil, mtproto.ErrUsernameInvalid
	}

	existed, err := c.svcCtx.Dao.UserClient.UserCheckChannelUsername(c.ctx, &userpb.TLUserCheckChannelUsername{
		ChannelId: in.GetChannel().GetChannelId(),
		Username:  username,
	})
	if err != nil {
		return nil, err
	}
	if existed.GetPredicateName() == userpb.Predicate_usernameExistedNotMe {
		return mtproto.BoolFalse, nil
	}
	return mtproto.BoolTrue, nil
}
