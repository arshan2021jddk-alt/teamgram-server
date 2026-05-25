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
	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/biz/message/message"
)

// MessageGetPinnedMessageIdList
// message.getPinnedMessageIdList user_id:long peer_type:int peer_id:long = Vector<int>;
func (c *MessageCore) MessageGetPinnedMessageIdList(in *message.TLMessageGetPinnedMessageIdList) (*message.Vector_Int, error) {
	dialogId := mtproto.MakeDialogId(in.UserId, in.PeerType, in.PeerId)
	idList, _ := c.svcCtx.Dao.MessagesDAO.SelectPinnedMessageIdList(c.ctx, in.UserId, dialogId.A, dialogId.B)

	return &message.Vector_Int{Datas: idList}, nil
}
