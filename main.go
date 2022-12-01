package main

import (
	"fmt"
	"log"
	"net/http"
)

func index(res http.ResponseWriter, req *http.Request){
	fmt.Print(res,"GET request")
}

func form(res http.ResponseWriter, req *http.Request){
	fmt.Print(res,"POST request")
}

func main() {
	fileServer:= http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/index", index)
	http.HandleFunc("/form", form)

	fmt.Println("Server is running on port 5500")
	if err:=http.ListenAndServe("127.0.0.1:5500",nil);err!=nil{
		log.Fatal(err)
	}
}