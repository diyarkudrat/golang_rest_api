package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	}

	var player Player
	err = json.Unmarshal(bodyBytes, &player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	player.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	p.Lock()
	p.store[player.FirstName] = player
	defer p.Unlock()

}

func newPlayerHandlers() *playerHandlers {
	return &playerHandlers{
		store: map[string]Player{
			"test": Player{
				FirstName: "Diyar",
				LastName:  "Kudrat",
				Position:  "OG",
				Height:    "5'11",
				Weight:    245,
				State:     "CA",
			},
		},
	}
}

func main() {
	playerHandlers := newPlayerHandlers()
	http.HandleFunc("/players", playerHandlers.players)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
