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

// SetupAPI - set global variables in api
func SetupAPI(pl *string, s map[string]db.User) {
	pukkalink = pl
	sessions = s
}

// SessionHandler - handle creation of sessions
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)

		w.Write(nil)
	} else if r.Method == http.MethodPost {
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

		uuid := sessionToken(carp.User)

		sessions[uuid] = carp.User

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(uuid))
	} else {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)
	}
}

func sessionToken(user db.User) string {
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)

		w.Write(nil)
	} else if r.Method == http.MethodPost {
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
	} else {
		w.WriteHeader(http.StatusUnauthorized)

		w.Write(nil)
	}
}
