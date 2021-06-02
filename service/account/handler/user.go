package handler

import (
	"context"
	"filestore-server/common"
	"filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/handler"
	"filestore-server/service/account/proto/proto"
	"filestore-server/util"
	"fmt"
)

type User struct{}
//Signup 处理用户注册请求
func (u *User)Signup(ctx context.Context, req *proto.ReqSignup,res *proto.RespSignup) error{
	username:=req.Username
	passwd:=req.Password
	fmt.Printf("user1:%s,pass2:%s",username,passwd)
	//存在一个之后 不在成功的问题 ，就是唯一索引值  没有插入，插入之后 就可以刷新
	if len(username)<3||len(passwd)<5{

		res.Code=-1
		res.Message="注册参数无效"
		return nil
	}
	enc_passwd:=util.Sha1([]byte(passwd+config.PasswordSalt))
	fmt.Println(enc_passwd,username)
	suc:=dblayer.UserSignUp(username,enc_passwd,"15210187668")
	if suc{
		res.Code=common.StatusOK
		res.Message="注册成功"
	}else{
		res.Code=common.StatusRegisterFailed
		res.Message="注册失败"
	}
	return  nil
}
func (u *User)Signin(ctx context.Context, req *proto.ReqSignin,res *proto.RespSignin) error{
	username:=req.Username
	password:=req.Password
	encPasswd:=util.Sha1([]byte(password+config.PasswordSalt))
	pwdChecked:=dblayer.UserSignin(username,encPasswd)
	if !pwdChecked{
		res.Message="登录失败"
		res.Code=-2
		return nil
	}
	//生成访问凭证
	token:=handler.GenToken(username)
	upRes:=dblayer.UpdateToken(username,token)
	if !upRes{

		res.Message="登录失败"
		res.Code=-2
		return nil
	}

	res.Code=0
	res.Message="ok"
	res.Token=token
	return nil
}

func (u *User)UserInfo(ctx context.Context, req *proto.ReqUserInfo,res *proto.RespUserInfo) error{
	username:=req.Username
	user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		return nil
	}

	res.LastActiveAt=user.LastActiveAt
	res.SignupAt=user.SignupAt
	return nil
}