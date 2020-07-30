package main

import "net/http"

type Player struct {
	FirstName string `json: "First Name"`
	LastName  string `json: "Last Name"`
	Height    string `json: "Height"`
	Weight    int    `json: "Weight"`
	State     string `json: "State"`
}

type playerHandlers struct {
	store map[string]Player
}

func (p *playerHandlers) get(w http.ResponseWriter, r *http.Request) {

}

func newPlayerHandlers() *playerHandlers {
	return &playerHandlers{
		store: map[string]Player{
			"test": Player{
				FirstName: "Diyar",
				LastName:  "Kudrat",
				Height:    "5'11",
				Weight:    245,
				State:     "CA",
			},
		},
	}
}

func main() {
	playerHandlers := newPlayerHandlers()
	http.HandleFunc("/players", playerHandlers.get)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
