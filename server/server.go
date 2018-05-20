package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type workRequest struct {
	event []byte
}

func (w *workRequest) DoWork() {
	fmt.Println("Event:", string(w.event))
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
		v, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Couldn't parse body", http.StatusBadRequest)
			return
		}
		go wp.Run(&workRequest{event: v})
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.Error(w, "POST method only", http.StatusMethodNotAllowed)
}
