package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Name string `json:"name"`
}

// The movieList is considered to be our Movies database
var moviesList []Movie

func main() {
	fmt.Println("=== MAIN START === MAIN START === MAIN START ===")

	// Call the function to initlize the Movies List (our movies database)
	initlizeMoviesList()

	// Setting-up our routes and their handlers
	r := mux.NewRouter()

	r.HandleFunc("/movie/create", createMovie).Methods("POST")
	r.HandleFunc("/movie/get", getMovie).Methods("GET")
	r.HandleFunc("/movie/update", updateMovie).Methods("POST")
	r.HandleFunc("/movie/delete", deleteMovie).Methods("POST")
	r.HandleFunc("/movies", getAllMovies).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("view").HTTPBox()))


	fmt.Println("===> STARTING SERVER ON PORT: 5000")

	// Operating the server with the created handler "r"
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal("Can not connect to the server", err)
	}
}

func initlizeMoviesList() {
	moviesList = append(moviesList, Movie{ID: "1", Title: "The Batman", Director: &Director{Name: "Rayan Noos"}})
	moviesList = append(moviesList, Movie{ID: "2", Title: "Man Of Steel", Director: &Director{Name: "DC Comics"}})
	moviesList = append(moviesList, Movie{ID: "3", Title: "Iron-man", Director: &Director{Name: "Marvel"}})
	moviesList = append(moviesList, Movie{ID: "4", Title: "Contarband", Director: &Director{Name: "Lee Apatche"}})
	moviesList = append(moviesList, Movie{ID: "5", Title: "The Dark Knight", Director: &Director{Name: "Rayan Noos"}})
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(moviesList)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-type", "application/json")

	r.ParseForm()

	for i, elem := range moviesList {
		if r.FormValue("id") == elem.ID {
			moviesList = append(moviesList[:i], moviesList[i+1:]...)
			json.NewEncoder(w).Encode(moviesList)
			return
		}
	}

}

func getMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-type", "application/json")

	r.ParseForm()

	for _, elem := range moviesList {
		if r.FormValue("id") == elem.ID {
			json.NewEncoder(w).Encode(elem)
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	r.ParseForm()
	newMovie := Movie{ID: r.FormValue("id"), Title: r.FormValue("title"), Director: &Director{Name: r.FormValue("dircName")}}

	// ID is created by getting the last added movie's ID and incremented by 1
	newId, _ := strconv.Atoi(moviesList[len(moviesList)-1].ID)
	newMovie.ID = strconv.Itoa(newId + 1)

	moviesList = append(moviesList, newMovie)

	json.NewEncoder(w).Encode(newMovie)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	r.ParseForm()

	for i, elem := range moviesList {
		if r.FormValue("id") == elem.ID {
			elem.Title = r.FormValue("title")
			elem.Director.Name = r.FormValue("dircName")
			moviesList[i] = elem
			json.NewEncoder(w).Encode(elem)
			return
		}
	}

}
