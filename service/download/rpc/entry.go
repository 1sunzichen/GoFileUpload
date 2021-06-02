package rpc

import (
	"context"
	dlProto "filestore-server/service/download/proto"
	cfg "filestore-server/service/download/config"
)

type Download struct {}
func (u *Download) DownloadEntry(	ctx context.Context,
	req *dlProto.ReqEntry,res *dlProto.RespEntry)error{

   res.Entry=cfg.DownloadEntry
   return nil
}
