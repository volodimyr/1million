package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

type WorkRequest struct {
	Token string `json:"token"`
}

func (w *WorkRequest) DoWork() {
	fmt.Println("Event:", w.Token)
}

func main() {
	wp, err := New(5096)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/perform", func(w http.ResponseWriter, req *http.Request) {
		performHandler(wp, w, req)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		wp.Shutdown()
		log.Fatalln(err)
		os.Exit(1)
	}
}

func performHandler(wp *WorkerPool, w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if req.Method == http.MethodPost {
		dec := json.NewDecoder(req.Body)
		var event WorkRequest
		err := dec.Decode(&event)
		if err != nil {
			http.Error(w, "Couldn't parse body", http.StatusBadRequest)
			return
		}
		go wp.Run(&event)
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.Error(w, "POST method only", http.StatusMethodNotAllowed)
}
