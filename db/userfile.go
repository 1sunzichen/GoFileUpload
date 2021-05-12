package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct{
	UserName string
	FileHash string
	FileName string
	UploadAt string
	LastUpdate string
	FileSize int
}
func OnUserFileUploadFinished(username,filehash,filename string,filesize int64)bool {
	stmt,err:=mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`) values(?,?,?,?,?) ")
	if err!=nil{
		fmt.Println(username,filehash,filename,filesize);
		return false
	}
	defer stmt.Close()
	_,err=stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err!=nil{
		fmt.Printf("errorTwo,%s %s %s %s",username,filehash,filename,filesize);
		return false
	}
	return true
}
//批量获取用户信息·
func QueryUserFileMatas(username string,limit int)([]UserFile,error){
	stmt,err:=mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?")
	if err!=nil{
		fmt.Printf("err=%v",err)
		return nil, err
	}
	defer stmt.Close()
	rows,err:=stmt.Query(username,limit)
	if err!=nil{
		fmt.Printf("err2=%v",err)
		return nil,err
	}
	var userFiles []UserFile
	for rows.Next(){
		ufile:=UserFile{}
		err=rows.Scan(&ufile.FileHash,&ufile.FileName,&ufile.FileSize,&ufile.UploadAt,
			&ufile.LastUpdate)
		if err!=nil{
			fmt.Println(err.Error())
			break
		}
		userFiles=append(userFiles,ufile)
	}
	return userFiles,nil

}
