package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/svasandani/terrabaq/internal/db"
	"github.com/svasandani/terrabaq/internal/util"
)

var pukkalink *string

var sessions map[string]db.User
var usrtouuid map[string]string

// Middleware - chain all middleware handlers in one nice convenient function :))
func Middleware(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return CorsHandler(PreflightRequestHandler(PostHandler(JSONHandler(fn))))
}

// CorsHandler - set all CORS headers
func CorsHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fn(w, r)
	})
}

// PostHandler - ensure all requests to API are posts
func PostHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fn(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)

			w.Write([]byte("Please use POST requests only."))
		}
	})
}

// JSONHandler - ensure all requests have JSON payloads
func JSONHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			fn(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)

			w.Write([]byte("Please submit JSON payloads only."))
		}
	})
}

// PreflightRequestHandler - respond with OK on CORS preflight check
func PreflightRequestHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			w.Write(nil)
		} else {
			fn(w, r)
		}
	})
}

// SetupAPI - set global variables in api
func SetupAPI(pl *string, s map[string]db.User, u map[string]string) {
	pukkalink = pl
	sessions = s
	usrtouuid = u
}

// SessionHandler - handle creation of sessions
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	var ok bool

	body, err := ioutil.ReadAll(r.Body)

	if ok = util.CheckHTTPError("error reading request body", err, w); !ok {
		return
	}

	car := db.ClientAccessRequest{}
	err = json.Unmarshal(body, &car)

	if ok = util.CheckHTTPError("error unmarshalling request body", err, w); !ok {
		return
	}

	carstring, err := json.Marshal(car)

	if ok = util.CheckHTTPError("error marshalling request body", err, w); !ok {
		return
	}

	link := *pukkalink + "client/auth"

	resp, err := http.Post(link, "application/json", bytes.NewBuffer(carstring))

	if ok = util.CheckHTTPError("error posting to pukka", err, w); !ok {
		return
	}

	defer resp.Body.Close()

	pukkaBody, err := ioutil.ReadAll(resp.Body)

	if ok = util.CheckHTTPError("error reading pukka response body", err, w); !ok {
		return
	}

	carp := db.ClientAccessResponse{}
	err = json.Unmarshal(pukkaBody, &carp)

	if ok = util.CheckHTTPError("error unmarshalling pukka response body", err, w); !ok {
		return
	}

	usrstring, err := json.Marshal(carp.User)

	if ok = util.CheckHTTPError("error marshalling user", err, w); !ok {
		return
	}

	uuid := sessionToken(string(usrstring))

	sessions[uuid] = carp.User
	usrtouuid[string(usrstring)] = uuid

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(uuid))
}

func sessionToken(user string) string {
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

// EnqueueHandler - handle enqueue requests
func EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	var ok bool

	body, err := ioutil.ReadAll(r.Body)

	if ok = util.CheckHTTPError("error reading request body", err, w); !ok {
		return
	}

	er := db.EnqueueRequest{}
	err = json.Unmarshal(body, &er)

	if ok = util.CheckHTTPError("error unmarshalling request body", err, w); !ok {
		return
	}

	user, ok := sessions[er.SessionToken]

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)

		return
	}

	userstring, err := json.Marshal(user)

	if ok = util.CheckHTTPError("error marshalling request body", err, w); !ok {
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(userstring)
}

// UpdateHandler - handle requests to update user roles
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var ok bool

	body, err := ioutil.ReadAll(r.Body)

	if ok = util.CheckHTTPError("error reading request body", err, w); !ok {
		return
	}

	ur := db.UpdateRequest{}
	err = json.Unmarshal(body, &ur)

	if ok = util.CheckHTTPError("error unmarshalling request body", err, w); !ok {
		return
	}

	usrstring, err := json.Marshal(ur.OldUser)

	if ok = util.CheckHTTPError("error marshalling user", err, w); !ok {
		return
	}

	uuid, ok := usrtouuid[string(usrstring)]

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)

		return
	}

	usrstring, err = json.Marshal(ur.NewUser)

	if ok = util.CheckHTTPError("error marshalling user", err, w); !ok {
		return
	}

	sessions[uuid] = ur.NewUser
	usrtouuid[string(usrstring)] = uuid

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(uuid))
}
