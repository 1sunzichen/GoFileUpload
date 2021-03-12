package meta
import(
mydb "filestore-server/db"
)
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}
var fileMetas map[string]FileMeta



func init(){
	fileMetas=make(map[string]FileMeta)
}
//新增/更新文件原信息
func UpdateFileMeta(fmeta FileMeta)  {
	fileMetas[fmeta.FileSha1]=fmeta
}
//新增/更新文件元信息 db
func UpdateFileMetaDb(fmeta FileMeta)bool{
	return mydb.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)

}
func GetFileMeta(fileSha1 string)FileMeta{
	return fileMetas[fileSha1]
}
//从 mysql 获取文件元信息
func GetFileMetaDB(fileSha1 string)(FileMeta,error){
	tfile,err:=mydb.GetFileMeta(fileSha1)
	if err!=nil{
		return FileMeta{},err
	}
	fmeta:=FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta,nil
}

//RemoveFileMeta 删除
func RemoveFileMeta(fileSha1 string){

	delete(fileMetas,fileSha1)
}

