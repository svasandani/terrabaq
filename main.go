package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var pukkalink *string

var sessions map[string]User

// Client - application requesting user data
type Client struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri"`
}

// User - export user struct for http
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}

// Role - export generic Role struct
type Role struct {
	Type       string `json:"type"`
	ResourceID string `json:"resource_id"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
	GrantType string `json:"grant_type"`
	AuthCode  string `json:"auth_code"`
	Client    Client `json:"client"`
}

// ClientAccessResponse - struct for responding to client access request
type ClientAccessResponse struct {
	User User `json:"user"`
}

type EnqueueRequest struct {
	SessionToken string `json:"session_token"`
}

type EnqueueResponse struct {
	User User `json:"user"`
}

func main() {
	pukkalink = flag.String("pukkaLink", "https://pukka.terraling.com/", "Link for email redirect")

	port := flag.String("port", "3000", "Port to serve Terrabaq")

	flag.Parse()

	sessions = make(map[string]User)

	http.HandleFunc("/", generateSessionToken)
	http.HandleFunc("/enqueue", enqueue)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func generateSessionToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)

		w.Write(nil)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println("error reading request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		car := ClientAccessRequest{}
		err = json.Unmarshal(body, &car)

		if err != nil {
			log.Println("error unmarshalling request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		carstring, err := json.Marshal(car)

		if err != nil {
			log.Println("error marshalling request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		link := *pukkalink + "client/auth"

		resp, err := http.Post(link, "application/json", bytes.NewBuffer(carstring))

		if err != nil {
			log.Println("error posting to pukka:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		defer resp.Body.Close()

		pukkaBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("error reading pukka response body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		carp := ClientAccessResponse{}
		err = json.Unmarshal(pukkaBody, &carp)

		if err != nil {
			log.Println("error unmarshalling pukka response body:")
			log.Println(string(pukkaBody))
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		uuid := sessionToken(carp.User)

		sessions[uuid] = carp.User

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(uuid))
	} else {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)
	}
}

func sessionToken(user User) string {
	randInt := rand.Int63n(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano() + randInt)

	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%x", b)

	return uuid
}

func enqueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)

		w.Write(nil)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println("error reading request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		er := EnqueueRequest{}
		err = json.Unmarshal(body, &er)

		if err != nil {
			log.Println("error unmarshalling request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		user, ok := sessions[er.SessionToken]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)

			w.Write(nil)

			return
		}

		userstring, err := json.Marshal(user)

		if err != nil {
			log.Println("error marshalling request body:")
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)

			w.Write(nil)

			return
		}

		w.WriteHeader(http.StatusOK)

		w.Write(userstring)
	} else {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)
	}
}
