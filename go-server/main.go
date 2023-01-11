package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(res http.ResponseWriter, req *http.Request){
	if err:= req.ParseForm(); err!=nil{
		fmt.Fprintf(res,"ParseForm error %v",err)
		return
	}
	fmt.Fprintf(res,"POST request successful")
	name:=req.FormValue("Name")
	address:=req.FormValue("Address")
	fmt.Fprintf(res,"\nName = %v",name)
	fmt.Fprintf(res,"\nAddress = %v",address)
}

func helloHandler(res http.ResponseWriter, req *http.Request){
	if req.URL.Path!="/hello"{
		http.Error(res,"404 not found",http.StatusNotFound)
		return
	}
	if req.Method!="GET"{
		http.Error(res,"404 not found",http.StatusNotFound)
		return
	}
	fmt.Fprintf(res,"hello")
}

func main() {
	fileserver := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileserver)
	http.HandleFunc("/form",formHandler)
	http.HandleFunc("/hello",helloHandler)

	fmt.Println("Server started at port 8080")

	if err:=http.ListenAndServe(":8080",nil);err!=nil{
		log.Fatal(err)
	}
}