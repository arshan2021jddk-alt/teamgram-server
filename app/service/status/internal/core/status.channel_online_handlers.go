/*
 * Created from 'scheme.tl' by 'mtprotoc'
 *
 * Copyright (c) 2021-present, Teamgram Studio (https://teamgram.io).
 * All rights reserved.
 */

package core

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/status/status"
)

const (
	channelUsersPrefix = "channel_online_users"
	userChannelsPrefix = "user_online_channels"
)

func getChannelUsersKey(channelID int64) string {
	return fmt.Sprintf("%s#%d", channelUsersPrefix, channelID)
}

func getUserChannelsKey(userID int64) string {
	return fmt.Sprintf("%s#%d", userChannelsPrefix, userID)
}

func uniquePositiveInt64(in []int64) []int64 {
	if len(in) == 0 {
		return nil
	}
	m := make(map[int64]struct{}, len(in))
	out := make([]int64, 0, len(in))
	for _, v := range in {
		if v <= 0 {
			continue
		}
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// StatusGetChannelOnlineUsers
// status.getChannelOnlineUsers channel_id:long = Vector<long>;
func (c *StatusCore) StatusGetChannelOnlineUsers(in *status.TLStatusGetChannelOnlineUsers) (*status.Vector_Long, error) {
	if in.GetChannelId() <= 0 {
		return nil, fmt.Errorf("status.getChannelOnlineUsers - invalid channelId: %d", in.GetChannelId())
	}

	rMap, err := c.svcCtx.Dao.KV.HgetallCtx(c.ctx, getChannelUsersKey(in.GetChannelId()))
	if err != nil {
		c.Logger.Errorf("status.getChannelOnlineUsers(channelId=%d) error: %v", in.GetChannelId(), err)
		return nil, err
	}

	onlineUsers := make([]int64, 0, len(rMap))
	for uid := range rMap {
		v, err := strconv.ParseInt(uid, 10, 64)
		if err != nil || v <= 0 {
			continue
		}
		onlineUsers = append(onlineUsers, v)
	}
	sort.Slice(onlineUsers, func(i, j int) bool { return onlineUsers[i] < onlineUsers[j] })

	return status.MakeTLVectorLong(&status.Vector_Long{Datas: onlineUsers}).To_Vector_Long(), nil
}

// StatusSetUserChannelsOnline
func (c *StatusCore) StatusSetUserChannelsOnline(in *status.TLStatusSetUserChannelsOnline) (*mtproto.Bool, error) {
	if in.GetUserId() <= 0 {
		return nil, fmt.Errorf("status.setUserChannelsOnline - invalid userId: %d", in.GetUserId())
	}

	uid := strconv.FormatInt(in.GetUserId(), 10)
	userChannelsKey := getUserChannelsKey(in.GetUserId())
	for _, ch := range uniquePositiveInt64(in.GetChannels()) {
		chKey := getChannelUsersKey(ch)
		err := c.svcCtx.Dao.KV.HsetCtx(c.ctx, chKey, uid, "1")
		if err != nil {
			return nil, err
		}
		err = c.svcCtx.Dao.KV.ExpireCtx(c.ctx, chKey, c.svcCtx.Config.StatusExpire)
		if err != nil {
			return nil, err
		}
		err = c.svcCtx.Dao.KV.HsetCtx(c.ctx, userChannelsKey, strconv.FormatInt(ch, 10), "1")
		if err != nil {
			return nil, err
		}
	}

	err := c.svcCtx.Dao.KV.ExpireCtx(c.ctx, userChannelsKey, c.svcCtx.Config.StatusExpire)
	if err != nil {
		return nil, err
	}
	return mtproto.BoolTrue, nil
}

// StatusSetUserChannelsOffline
func (c *StatusCore) StatusSetUserChannelsOffline(in *status.TLStatusSetUserChannelsOffline) (*mtproto.Bool, error) {
	if in.GetUserId() <= 0 {
		return nil, fmt.Errorf("status.setUserChannelsOffline - invalid userId: %d", in.GetUserId())
	}

	uid := strconv.FormatInt(in.GetUserId(), 10)
	userChannelsKey := getUserChannelsKey(in.GetUserId())
	for _, ch := range uniquePositiveInt64(in.GetChannels()) {
		err := c.svcCtx.Dao.KV.HdelCtx(c.ctx, getChannelUsersKey(ch), uid)
		if err != nil {
			return nil, err
		}
		err = c.svcCtx.Dao.KV.HdelCtx(c.ctx, userChannelsKey, strconv.FormatInt(ch, 10))
		if err != nil {
			return nil, err
		}
	}

	return mtproto.BoolTrue, nil
}

// StatusSetChannelUserOffline
func (c *StatusCore) StatusSetChannelUserOffline(in *status.TLStatusSetChannelUserOffline) (*mtproto.Bool, error) {
	if in.GetChannelId() <= 0 || in.GetUserId() <= 0 {
		return nil, fmt.Errorf("status.setChannelUserOffline - invalid params: channelId=%d userId=%d", in.GetChannelId(), in.GetUserId())
	}

	err := c.svcCtx.Dao.KV.HdelCtx(c.ctx, getChannelUsersKey(in.GetChannelId()), strconv.FormatInt(in.GetUserId(), 10))
	if err != nil {
		return nil, err
	}
	err = c.svcCtx.Dao.KV.HdelCtx(c.ctx, getUserChannelsKey(in.GetUserId()), strconv.FormatInt(in.GetChannelId(), 10))
	if err != nil {
		return nil, err
	}
	return mtproto.BoolTrue, nil
}

// StatusSetChannelUsersOnline
func (c *StatusCore) StatusSetChannelUsersOnline(in *status.TLStatusSetChannelUsersOnline) (*mtproto.Bool, error) {
	if in.GetChannelId() <= 0 {
		return nil, fmt.Errorf("status.setChannelUsersOnline - invalid channelId: %d", in.GetChannelId())
	}

	for _, uid := range uniquePositiveInt64(in.GetId()) {
		err := c.svcCtx.Dao.KV.HsetCtx(c.ctx, getChannelUsersKey(in.GetChannelId()), strconv.FormatInt(uid, 10), "1")
		if err != nil {
			return nil, err
		}
		err = c.svcCtx.Dao.KV.HsetCtx(c.ctx, getUserChannelsKey(uid), strconv.FormatInt(in.GetChannelId(), 10), "1")
		if err != nil {
			return nil, err
		}
		err = c.svcCtx.Dao.KV.ExpireCtx(c.ctx, getUserChannelsKey(uid), c.svcCtx.Config.StatusExpire)
		if err != nil {
			return nil, err
		}
	}
	err := c.svcCtx.Dao.KV.ExpireCtx(c.ctx, getChannelUsersKey(in.GetChannelId()), c.svcCtx.Config.StatusExpire)
	if err != nil {
		return nil, err
	}
	return mtproto.BoolTrue, nil
}

// StatusSetChannelOffline
func (c *StatusCore) StatusSetChannelOffline(in *status.TLStatusSetChannelOffline) (*mtproto.Bool, error) {
	if in.GetChannelId() <= 0 {
		return nil, fmt.Errorf("status.setChannelOffline - invalid channelId: %d", in.GetChannelId())
	}

	rMap, err := c.svcCtx.Dao.KV.HgetallCtx(c.ctx, getChannelUsersKey(in.GetChannelId()))
	if err != nil {
		return nil, err
	}
	for uid := range rMap {
		uid64, err := strconv.ParseInt(uid, 10, 64)
		if err != nil || uid64 <= 0 {
			continue
		}
		err = c.svcCtx.Dao.KV.HdelCtx(c.ctx, getUserChannelsKey(uid64), strconv.FormatInt(in.GetChannelId(), 10))
		if err != nil {
			return nil, err
		}
	}
	err = c.svcCtx.Dao.KV.DelCtx(c.ctx, getChannelUsersKey(in.GetChannelId()))
	if err != nil {
		return nil, err
	}
	return mtproto.BoolTrue, nil
}
