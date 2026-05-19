package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt="*#890"
)
// Handle user signup request
func SignupHandler(w http.ResponseWriter,r *http.Request){
	if r.Method==http.MethodGet{
		data,err:=ioutil.ReadFile("./static/view/signup.html")
		if err!=nil{

		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username:=r.Form.Get("username")
	passwd:=r.Form.Get("password")
	// There was an issue where it would not succeed after the first one;
	// the unique index value was not inserted, and after insertion it can be refreshed.
	phone:=r.Form.Get("phone")
	if len(username)<3||len(passwd)<5{
		w.Write([]byte("Invaild parameter"))
		return
	}
	enc_passwd:=util.Sha1([]byte(passwd+pwd_salt))
	fmt.Println(enc_passwd,username)
	suc:=dblayer.UserSignUp(username,enc_passwd,phone)
	if suc{
		w.Write([]byte("success"))
	}else{
		w.Write([]byte("Fail"))
	}
}
type result struct {
	url string
}
// SignInHandler login endpoint
func SignInHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	// Verify username and password
	username:=r.Form.Get("username")
	password:=r.Form.Get("password")
	encPasswd:=util.Sha1([]byte(password+pwd_salt))
	pwdChecked:=dblayer.UserSignin(username,encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		return
	}
	// Generate access token
	token:=GenToken(username)
	upRes:=dblayer.UpdateToken(username,token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}
	// Redirect to home page after successful login
	//url:=result{url:"http://"+r.Host+"/static/view/home.html"}
	//http.Redirect(w,r,"http://"+r.Host+"/static/view/home.html",http.StatusFound)
	//w.Write([]byte(`{"code":401,"msg":"http://`+r.Host+`/static/view/home.html"}`))
		resp:=util.RespMsg{
			Code:0,
			Msg:"ok",
			Data:struct{
				Location string
				Username string
				Token string
			}{
				Location:"http://"+r.Host+"/static/view/home.html",
				Username:username,
				Token:token,
			},
		}
		w.Write(resp.JSONBytes())
}
func UserInfoHandler(w http.ResponseWriter,r *http.Request){

	// 1. Parse request parameters
	r.ParseForm()
	username:=r.Form.Get("username")
	//token:=r.Form.Get("token")
	//// 2. Verify whether the token is valid
	//isVaildToken:=IsTokenVaild(token)
	//if !isVaildToken{
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	// 3. Query user info
	user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4. Assemble and respond with user data
	resp:=util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenVaild(token string)bool{
	// TODO Check token validity period
	if len(token)!=40{
		return false
	}
	// From the tbl_user table
	return true
}
func GenToken(username string)string{
	// 40-character
	ts:=fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix:=util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix+ts[:8]
}
