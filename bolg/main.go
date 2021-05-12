package main
// import(
// 	"github.com/scottkiss/gosk"
	
// )
import (
	// "filestore-server/handler"
	"fmt"
	"net/http"
)
func main(){
	// gosk.Build();
    // gosk.Run(":8037")
	fs := http.FileServer(http.Dir("dist/"))
	http.Handle("/", http.StripPrefix("/", fs))
	err:=http.ListenAndServe(":3000",nil)
	if err!=nil{
		fmt.Printf("Filed to start server",err.Error())
	}
}
