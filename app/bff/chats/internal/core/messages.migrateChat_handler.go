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
)

// MessagesMigrateChat
// messages.migrateChat#a2875319 chat_id:long = Updates;
func (c *ChatsCore) MessagesMigrateChat(in *mtproto.TLMessagesMigrateChat) (*mtproto.Updates, error) {
	if in.GetChatId() <= 0 {
		err := mtproto.ErrChatIdInvalid
		c.Logger.Errorf("messages.migrateChat - invalid chat_id: %v", err)
		return nil, err
	}

	// TODO: wire to channel service once channel creation/migration backend is available.
	c.Logger.Errorf("messages.migrateChat - error: backend channel migration flow is not implemented")
	return nil, mtproto.ErrMethodNotImpl
}
