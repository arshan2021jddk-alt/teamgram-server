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
	"github.com/teamgram/teamgram-server/app/service/dfs/dfs"
	"github.com/teamgram/teamgram-server/app/service/media/media"
)

// MediaUploadStickerFile
// media.uploadStickerFile flags:# owner_id:long file:InputFile thumb:flags.0?InputFile mime_type:string file_name:string document_attribute_sticker:DocumentAttribute = Document;
func (c *MediaCore) MediaUploadStickerFile(in *media.TLMediaUploadStickerFile) (*mtproto.Document, error) {
	if in.GetFile() == nil {
		return nil, mtproto.ErrMediaInvalid
	}

	attrs := make([]*mtproto.DocumentAttribute, 0, 1)
	if in.GetDocumentAttributeSticker() != nil {
		attrs = append(attrs, in.GetDocumentAttributeSticker())
	}

	inputMedia := mtproto.MakeTLInputMediaUploadedDocument(&mtproto.InputMedia{
		File:       in.GetFile(),
		Thumb:      in.GetThumb(),
		MimeType:   in.GetMimeType(),
		Attributes: attrs,
	}).To_InputMedia()

	document, err := c.svcCtx.Dao.DfsClient.DfsUploadDocumentFileV2(c.ctx, &dfs.TLDfsUploadDocumentFileV2{
		Creator: in.GetOwnerId(),
		Media:   inputMedia,
	})
	if err != nil {
		c.Logger.Errorf("media.uploadStickerFile - error: %v", err)
		return nil, err
	}

	if len(document.GetThumbs()) > 0 {
		if err = c.svcCtx.Dao.SavePhotoSizeV2(c.ctx, document.GetId(), document.GetThumbs()); err != nil {
			c.Logger.Errorf("media.uploadStickerFile - save thumbs error: %v", err)
			return nil, err
		}
	}
	c.svcCtx.Dao.SaveDocumentV2(c.ctx, in.GetFileName(), document)

	return document, nil
}
