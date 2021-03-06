package coaster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Coaster struct
type Coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           string `json:"id"`
	InPark       string `json:"inPark"`
	Height       int    `json:"height"`
}

// Store exposes an in-memory MAP of all the coasters
// Implements the ServeHTTP Handler interface
type Store struct {
	mutex sync.Mutex
	Store map[string]Coaster
}

func (s *Store) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case "GET":
		// a switch in a switch...
		switch len(p) {
		case 2:
			s.get(w, r)
			return
		case 3:
			if p[2] == "random" {
				s.getRandomCoaster(w, r)
			} else {
				s.getCoaster(w, r)
			}
			return
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Wrong Path"))
			return
		}
	case "POST":
		s.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
}

func (s *Store) getRandomCoaster(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(s.Store))

	// From Store
	s.mutex.Lock()
	i := 0
	for id := range s.Store {
		ids[i] = id
		i++
	}
	s.mutex.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Coasters in Store"))
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		rng := rand.Intn(len(ids))
		target = ids[rng]
	}

	http.Redirect(w, r, "/coaster/"+target, http.StatusSeeOther)
}

func (s *Store) getCoaster(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")

	// From Store
	s.mutex.Lock()
	cstr, ok := s.Store[p[2]]
	s.mutex.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Couldnt find Coaster id: %s", p[2])))
		return
	}

	json, err := json.Marshal(cstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *Store) get(w http.ResponseWriter, r *http.Request) {
	cs := make([]Coaster, len(s.Store))

	// From Store
	s.mutex.Lock()
	i := 0
	for _, c := range s.Store {
		cs[i] = c
		i++
	}
	s.mutex.Unlock()

	json, err := json.Marshal(cs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *Store) post(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var c Coaster
	err = json.Unmarshal(b, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type: 'application/json' but got %s", ct)))
		return
	}

	// set a unique ID
	c.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// From Store
	s.mutex.Lock()
	s.Store[c.ID] = c
	s.mutex.Unlock()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
