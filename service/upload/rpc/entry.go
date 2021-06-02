package rpc

import (
	"context"
	"filestore-server/service/upload/config"
	upProto "filestore-server/service/upload/proto"
)
type Upload struct{}
func (u *Upload)UploadEntry(ctx context.Context, req *upProto.ReqEntry, res *upProto.RespEntry) error{
	res.Entry=config.UploadEntry
	return nil
}

