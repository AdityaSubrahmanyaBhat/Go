package main

import(
	"fmt"
	"log"
	"net/http"
)

func helloHandler(resw http.ResponseWriter,req *http.Request){
	if(req.URL.Path!="/hello"){
		http.Error(resw,"404 not found",http.StatusNotFound)
		return
	}
	if(req.Method!="GET"){
		http.Error(resw,"Illegal method",http.StatusNotFound)
		return
	}
	fmt.Fprintf(resw,"Hello there , this is what was supposed to be printed")
}

func formHandler(resw http.ResponseWriter,req *http.Request){
	if(req.URL.Path!="/form"){
		http.Error(resw,"404 not found",http.StatusNotFound)
		return
	}
	name:=req.FormValue("name")
	specialization:=req.FormValue("specialization")
	fmt.Fprintf(resw,"Name = %s\nSpecialization = %s",name,specialization)
}

func main(){
	fmt.Println("Server started at port 8080")
	fileServer:=http.FileServer(http.Dir("./static"))
	http.Handle("/",fileServer)
	http.HandleFunc("/hello",helloHandler)
	http.HandleFunc("/form",formHandler)

	err:=http.ListenAndServe(":8080",nil)
	if(err!=nil){
		log.Fatal(err)
	}
}

