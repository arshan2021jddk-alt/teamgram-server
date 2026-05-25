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

// MessageGetPeerUserMessageId
// message.getPeerUserMessageId user_id:long peer_user_id:long msg_id:int = Int32;
func (c *MessageCore) MessageGetPeerUserMessageId(in *message.TLMessageGetPeerUserMessageId) (*mtproto.Int32, error) {
	myDO, err := c.svcCtx.Dao.MessagesDAO.SelectByMessageId(c.ctx, in.UserId, in.MsgId)
	if err != nil || myDO == nil {
		return mtproto.MakeTLInt32(&mtproto.Int32{V: 0}).To_Int32(), nil
	}

	peerDialogId := mtproto.MakeDialogId(in.PeerUserId, mtproto.PEER_USER, in.UserId)
	peerList, _ := c.svcCtx.Dao.MessagesDAO.SelectDialogMessageIdList(c.ctx, in.PeerUserId, peerDialogId.A, peerDialogId.B)
	var peerMsgId int32
	for i := range peerList {
		if peerList[i].DialogMessageId == myDO.DialogMessageId {
			peerMsgId = peerList[i].UserMessageBoxId
			break
		}
	}

	return mtproto.MakeTLInt32(&mtproto.Int32{V: peerMsgId}).To_Int32(), nil
}
