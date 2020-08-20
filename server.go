package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Player struct {
	ID        string `json:"id"`
	FirstName string `json: "First Name"`
	LastName  string `json: "Last Name"`
	Position  string `json: "Position"`
	Height    string `json: "Height"`
	Weight    int    `json: "Weight"`
	State     string `json: "State"`
}

type playerHandlers struct {
	sync.Mutex
	store map[string]Player
}

func (p *playerHandlers) players(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p.get(w, r)
		return
	case "POST":
		p.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

}

func (p *playerHandlers) get(w http.ResponseWriter, r *http.Request) {
	players := make([]Player, len(p.store))

	p.Lock()
	// Lock so only one goroutine at a time can access the map
	i := 0
	for _, player := range p.store {
		players[i] = player
		i++
	}
	p.Unlock()

	jsonBytes, err := json.Marshal(players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (p *playerHandlers) getRandomPlayer(w http.ResponseWriter, r *http.Request) {
	idList := make([]string, len(p.store))
	p.Lock()
	i := 0
	for id := range p.store {
		idList[i] = id
		i++
	}
	defer p.Unlock()

	var target string
	if len(idList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(idList) == 1 {
		target = idList[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = idList[rand.Intn(len(idList))]
	}

	w.Header().Add("location", fmt.Sprintf("/players/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (p *playerHandlers) getPlayer(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		p.getRandomPlayer(w, r)
		return
	}

	p.Lock()
	player, ok := p.store[parts[2]]
	p.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (p *playerHandlers) post(w http.ResponseWriter, r *http.Request) {
	// json data
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contype := r.Header.Get("content-type")
	if contype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", contype)))
		return
	}

	var player Player
	err = json.Unmarshal(bodyBytes, &player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	player.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	p.Lock()
	p.store[player.ID] = player
	defer p.Unlock()

}

func newPlayerHandlers() *playerHandlers {
	return &playerHandlers{
		store: map[string]Player{},
	}
}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required ENV variable ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
		return
	}

	w.Write([]byte("<html><h1>Secret admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	playerHandlers := newPlayerHandlers()
	http.HandleFunc("/players", playerHandlers.players)
	http.HandleFunc("/players/", playerHandlers.getPlayer)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
