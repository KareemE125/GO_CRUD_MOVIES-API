package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	rice "github.com/GeertJohan/go.rice"
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

	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("view").HTTPBox()))
	//r.Handle("/", http.FileServer(http.Dir("./view")))

	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movies", getAllMovies).Methods("GET")

	fmt.Println("===> STARTING SERVER ON PORT: 3000")

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

	params := mux.Vars(r)

	for i, elem := range moviesList {
		if params["id"] == elem.ID {
			moviesList = append(moviesList[:i], moviesList[i+1:]...)
			json.NewEncoder(w).Encode(moviesList)
			return
		}
	}

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	params := mux.Vars(r)

	for _, elem := range moviesList {
		if params["id"] == elem.ID {
			json.NewEncoder(w).Encode(elem)
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var newMovie Movie
	json.NewDecoder(r.Body).Decode(&newMovie)

	// ID is created by getting the last added movie's ID and incremented by 1
	newId, _ := strconv.Atoi(moviesList[len(moviesList)-1].ID)
	newMovie.ID = strconv.Itoa(newId + 1)

	moviesList = append(moviesList, newMovie)

	json.NewEncoder(w).Encode(newMovie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)
	for i, elem := range moviesList {
		if params["id"] == elem.ID {
			json.NewDecoder(r.Body).Decode(&elem)

			// Replace the sent id by the real one
			// as we only want to update its data with the id remianing the same
			elem.ID = params["id"]
			moviesList[i] = elem
			json.NewEncoder(w).Encode(elem)
			return
		}
	}

}
