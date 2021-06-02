package handler

import (
	"context"
	"encoding/json"
	"filestore-server/common"
	"filestore-server/service/account/proto/proto"
	dblayer "filestore-server/db"
)

// UserFiles : 获取用户文件列表 QueryUserFileMatas
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dblayer.QueryUserFileMatas(req.Username, int(req.Limit))
	if err != nil {
		res.Code = common.StatusServerError
		return err
	}

	data, err := json.Marshal(dbResp)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

// UserFiles : 用户文件重命名
