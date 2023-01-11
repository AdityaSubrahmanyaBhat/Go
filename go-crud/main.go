package main

import (
	"encoding/json"
	"fmt"
	// f "go-crud/functions"
	// d "go-crud/modules/director"
	// m "go-crud/modules/movies"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct{
	ID string `json:"ID"`
	ISBN string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct{
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func getMovies(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type","application/json")
	json.NewEncoder(res).Encode(movies)
}

func getMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type","application/json")
	params:=mux.Vars(req)
	for _, item := range movies{
		if item.ID==params["ID"]{
			json.NewEncoder(res).Encode(item)
			return
		} 
	}
	json.NewEncoder(res).Encode(movies)
}

func deleteMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type","application/json")
	params:=mux.Vars(req)
	for index, item := range movies{
		if item.ID==params["ID"]{
			movies=append(movies[:index],movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(res).Encode(movies)
}

func createMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type","application/json")
	var movie Movie
	_=json.NewDecoder(req.Body).Decode(&movie)
	movie.ID=strconv.Itoa(rand.Intn(1000000000000))
	movies=append(movies, movie)
	json.NewEncoder(res).Encode(movie)
	
}

func updateMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type","application/json")
	params:=mux.Vars(req)
	for index, item := range movies{
		if item.ID==params["ID"]{
			movies=append(movies[:index],movies[index+1:]...)
			var movie Movie
			_=json.NewDecoder(req.Body).Decode(&movie)
			movie.ID=params["ID"]
			movies=append(movies,movie)
			json.NewEncoder(res).Encode(movie)
		}
	}
	json.NewEncoder(res).Encode(movies)
}


var movies []Movie

func main(){

	movies = append(movies, Movie{ID:"0",Title:"the matrix",ISBN:"123456",Director:&Director{FirstName:"aditya",LastName:"bhat"}})

	router:=mux.NewRouter()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies",createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}",updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}",deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8080...")
	// if err := http.ListenAndServe(":8080",nil);err!=nil{
	// 	log.Fatal("there is something wrong with us ",err)
	// }
	log.Fatal(http.ListenAndServe(":8080",router))
}