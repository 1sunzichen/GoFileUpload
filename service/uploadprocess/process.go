package  main

import (
	"filestore-server/route"
	//"fmt"
	//"net/http"

)

func main(){
	router:=route.Route()
	router.Run(":8080")

}
