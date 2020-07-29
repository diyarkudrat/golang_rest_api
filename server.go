package main

import "net/http"

func playersHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandlerFunc("/", "")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
