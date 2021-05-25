package main

import lua "github.com/yuin/gopher-lua"

func main() {
	l:=lua.NewState()
	defer l.Close()
	if err:=l.DoString(`print("hello world")`);err!=nil{
		panic(err)
	}

}
