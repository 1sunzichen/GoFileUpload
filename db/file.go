package db
import (
	mydb "filestore-server/db/mysql"
	"fmt"
)
//文件上传 接口
func OnFileUploadFinished(filehash string,filename string,filesize int64,fileaddr string)bool {
	stmt,err:=mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,"+
			"`file_addr`,`status` values(?,?,?,?,1) ",
		)
	if err!=nil{
		fmt.Println("Faile to prepare statementm err:"+err.Error())
		return false
	}
	defer stmt.Close()
	ret,err:=stmt.Exec(filehash,filename,filesize,fileaddr)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	if rf,err:=ret.RowsAffected();nil==err{
		if rf<=0{
			fmt.Printf("File with hash%s hash been upload before",filehash)
		}
		return true
	}
	return  false

}
