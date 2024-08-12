package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// this creates a blueprint for the movie and director information. we include
// a JSON tag to specify how the field will be represented when it is marshaled
// to JSON. each movie also aggregates a director - it is associated with
// an existing director instance, exemplified by using a pointer to a director.
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// creating a slice to store our movies in
var movies []Movie

// this function is a handler for GET requests to the /movies endpoint of the
// server. it takes in w, a ResponseWriter, which allows us to directly interact
// with the HTTP response, and it also takes in r, a pointer to a request.
func getMovies(w http.ResponseWriter, r *http.Request) {

	// this sets the "Content-Type" header of the HTTP response to JSON format.
	w.Header().Set("Content-Type", "application/json")

	// json.NewEncoder is an object that writes data that is in JSON format to w.
	// .Encode(movies) actually marshals the movies slice into JSON format and
	// then writes it to w.
	json.NewEncoder(w).Encode(movies)
}

// this function is a handler for DELETE requests to the /movies/{id} endpoint
// of the server.
func deleteMovie(w http.ResponseWriter, r *http.Request) {

	// this sets the "Content-Type" header of the HTTP response to JSON format.
	// NOTE: although this line is not particularly useful since we do not
	// actually write anything to the HTTP response in this function, leaving
	// it here is good practice and maintains code consistency.
	w.Header().Set("Content-Type", "application/json")

	// mux.Vars(r) takes in a Request and returns any URL variables in the
	// route pattern as a map. for this specific request, we would extract
	// whatever the client put in for {id} in the "/movies/{id}" route pattern.
	params := mux.Vars(r)

	// iterate through the slice of movies. if the movie at `i`'s ID has the
	// same contents as the id key in params, we will remove that movie from
	// `movies` and break from the for loop.
	for i := 0; i < len(movies); i++ {
		if movies[i].ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}
}

// this function is a handler for GET requests to the /movies/{id} endpoint
// of the server.
func getMovie(w http.ResponseWriter, r *http.Request) {

	// this sets the "Content-Type" header of the HTTP response to JSON format.
	w.Header().Set("Content-Type", "application/json")

	// mux.Vars(r) takes in a Request and returns any URL variables in the
	// route pattern as a map. for this specific request, we would extract
	// whatever the client put in for {id} in the "/movies/{id}" route pattern.
	params := mux.Vars(r)

	// iterate through the slice of movies. if the movie at `i`'s ID has the
	// same contents as the id key in params, we will write that movie in JSON
	// format to the HTTP response and return.
	for i := 0; i < len(movies); i++ {
		if movies[i].ID == params["id"] {
			json.NewEncoder(w).Encode(movies[i])
			return
		}
	}
}

// this function is a handler for POST requests to the /movies endpoint of the
// server. it will create a new movie and then send it back as a JSON response.
func createMovie(w http.ResponseWriter, r *http.Request) {

	// this sets the "Content-Type" header of the HTTP response to JSON format.
	w.Header().Set("Content-Type", "application/json")

	// declare a movie variable. then we unmarshal the movie data from the
	// request's body and store it in the value pointed to by `movie`
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)

	// generate a random integer from 0 - 999999, convert it to a string, and
	// set the movie's ID to it.
	movie.ID = strconv.Itoa(rand.Intn(1000000))

	// update `movies` to include the newly created movie
	movies = append(movies, movie)

	// marshal the movie back into JSON and write it to the HTTP response
	json.NewEncoder(w).Encode(movie)
}

// this function is a handler for PUT requests to the /movies/{id} endpoint of
// the server. it allows us to update the contents of a movie.
func updateMovie(w http.ResponseWriter, r *http.Request) {

	// this sets the "Content-Type" header of the HTTP response to JSON format.
	w.Header().Set("Content-Type", "application/json")

	// store the movie id specified in the route pattern as a key value pair
	params := mux.Vars(r)

	// iterate through `movies`. if a movie ID matches the ID extracted from
	// the route pattern, remove the existing version of that movie from
	// `movies`. then create a new movie with the updated information (see
	// createMovie) and append it to `movies`. write the updated movie in JSON
	// format to w and return.
	for i := 0; i < len(movies); i++ {
		if movies[i].ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			var movie Movie
			json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	// creating a mux.Router instance. this Router will allow us to create
	// Routes that match HTTP requests to the correct handler functions based on
	// the URL path that the request is made to.
	r := mux.NewRouter()

	// adding some movies to our slice so that when we send a GET request to
	// /movies, there will actually be a result
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "123456",
		Title: "Movie One",
		Director: &Director{
			Firstname: "Lebron",
			Lastname:  "James",
		},
	})
	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "654321",
		Title: "Movie Two",
		Director: &Director{
			Firstname: "Joe",
			Lastname:  "Biden",
		},
	})

	// actually defining the routes as stated above. for example, the line below
	// creates a route that handles GET requests to /movies by calling getMovies
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// starting the server and telling it to listen on port 8000 while using r
	// (the router we defined earlier) to handle any requests. if a non-nil
	// error is returned by ListenAndServe, we will log it and exit the program.
	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
