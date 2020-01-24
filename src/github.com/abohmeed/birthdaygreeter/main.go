package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var redisHostWrite string
var redisHostRead string
var redisPort string
var redisPassword string

const layoutISO = "2006-01-02"

//PostData a struct for holding the the data sent in the request
type PostData struct {
	Birthdate string `json:"dateOfBirth"`
}

func newServer() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(commonMiddleware)
	r.HandleFunc("/hello/{username}", handlePostBirthdate).Methods("POST")
	return r
}

func main() {
	var router = newServer()
	log.Println("Server starting on port 8080")
	log.Fatal("Application is running", http.ListenAndServe(":8080", router))
}

func handlePostBirthdate(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading the request body:", err)
	}
	routeVars := mux.Vars(r)
	var pd PostData
	if err = json.Unmarshal(b, &pd); err != nil {
		respondWithError(w, "Invalid request", 422)
		return
	}
	// This is useful to know what the users are sending (someone is trying to intentionally send incorrect data?)
	log.Print("Received user " + routeVars["username"] + " with birthdate " + pd.Birthdate)
	if t, err := time.Parse(layoutISO, pd.Birthdate); err != nil {
		respondWithError(w, `Invalid message format, should be "dateOfBirth":"yyyy-mm-dd"`, 422)
	} else {
		vars := mux.Vars(r)
		username := vars["username"]
		// We don't need the time part of the birthdate so let's remove it
		bd := t.Format("2006-01-02")
		daystillbdate := getTimeTillBirthdate(bd)
		if daystillbdate <= 0 && daystillbdate > -1 {
			json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username + "! Happy birthday"})
		} else {
			// The user's birthday is neither today nor 5 days ahead, so we just greet them
			json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username})
		}
	}
}

func respondWithError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func getTimeTillBirthdate(t string) int16 {
	// We assume that the date values stored in Redis are correct because we try to parse them before storing them
	bdate, _ := time.Parse(layoutISO, t)
	cdate := time.Now()
	// We need the month and day only of both dates
	cdate = time.Date(cdate.Year(), cdate.Month(), cdate.Day(), 0, 0, 0, 0, time.UTC)
	// Chnage the year part of the birthdate to be the current year
	bdate = time.Date(cdate.Year(), bdate.Month(), bdate.Day(), 0, 0, 0, 0, time.UTC)
	diff := bdate.Sub(cdate).Hours()
	return int16(diff / float64(24))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
