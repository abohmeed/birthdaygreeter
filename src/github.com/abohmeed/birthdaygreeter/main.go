package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var isHealthy bool = true

var redisHost string
var redisPort string
var appPort string

const layoutISO = "2006-01-02"

//PostData a struct for holding the the data sent in the request
type PostData struct {
	Birthdate string `json:"dateOfBirth"`
}

func setEnv() {
	if appPort = os.Getenv("APP_PORT"); appPort == "" {
		appPort = "3000"
	}
	if redisHost = os.Getenv("REDIS_HOST"); redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort = os.Getenv("REDIS_PORT"); redisPort == "" {
		redisPort = "6379"
	}
}

func newServer() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	r.HandleFunc("/hello/{username}", handleUpdateBirthdate).Methods("PUT")
	r.HandleFunc("/hello/{username}", handleQueryBirthdate).Methods("GET")
	return r
}

func main() {
	setEnv()
	// Just checking that Redis is reachable
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()
	ping(conn)
	var router = newServer()
	log.Println("Server starting on port", appPort)
	log.Fatal("Application is running", http.ListenAndServe(":"+appPort, router))
}

func handleUpdateBirthdate(w http.ResponseWriter, r *http.Request) {
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
		// fmt.Println(t)
		vars := mux.Vars(r)
		username := vars["username"]
		pool := newPool()
		conn := pool.Get()
		defer conn.Close()
		// We don't need the time part of the birthdate so let's remove it
		err = set(conn, username, t.Format("2006-01-02"))
		if err != nil {
			log.Println("Error while saving record to Redis:", err)
		}
		w.WriteHeader(204)
	}
}

func handleQueryBirthdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()
	if bd, err := get(conn, username); err != nil {
		// This is a new user who hasn't got any data yet
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username})
	} else {
		daystillbdate := getTimeTillBirthdate(bd)
		if daystillbdate <= 5 && daystillbdate > 4 {
			json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username + "! Your birthday is in 5 days"})
		} else if daystillbdate <= 0 && daystillbdate > -1 {
			json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username + "! Happy birthday"})
		} else {
			// The user's birthday is neither today nor 5 days ahead, so we just greet them
			json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username})
		}
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if isHealthy {
		json.NewEncoder(w).Encode("OK")
	} else {
		json.NewEncoder(w).Encode("Application failed")
	}
}
func respondWithError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func newPool() *redis.Pool {
	// We need to set the Redis connection settings for testing functions individually (not passing through main() function)
	setEnv()
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost+":"+redisPort)
			if err != nil {
				log.Println("Could not reach Redis", err)
				isHealthy = false
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	_, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}
	return nil
}
func set(c redis.Conn, key string, value string) error {
	_, err := c.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

// get executes the redis GET command
func get(c redis.Conn, key string) (string, error) {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	return s, nil
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
