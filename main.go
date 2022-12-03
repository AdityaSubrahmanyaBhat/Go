package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	ISBN string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

var movies []Movie

func getMovies(res http.ResponseWriter, req *http.Request ) {
	res.Header().Set("Content-Type","application/json")
	json.NewEncoder(res).Encode(movies)
}

func main() {

	movies=append(movies, Movie{ID: "1", ISBN: "12345", Title: "Movie One", Director: &Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", ISBN: "12346", Title: "Movie Two", Director: &Director{FirstName: "Steve", LastName: "Smith"}})
	movies=append(movies, Movie{ID:"3",ISBN:"12347",Title:"Movie Three",Director:&Director{FirstName:"Mary",LastName:"Jane"}})

	routes:=mux.NewRouter()
	routes.HandleFunc("/movies",getMovies).Methods("GET")
	routes.HandleFunc("/movies/{id}",getMovie).Methods("GET")
	routes.HandleFunc("movies",createMovie).Methods("POST")

	http.ListenAndServe("127.0.0.1/5500",routes)
	fmt.Println("Server started at port 5500")
}