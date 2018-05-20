package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	URL         = "http://localhost:8080/perform"
	contentType = "application/json"
)

var (
	reqPerM  uint64
	routines = 1
)

const (
	TIMEOUT = time.Second * 5
)

type Result struct {
	//failed sending request
	Failed bool
	//response status
	StatusCode    int
	LatencyMillis float64
}

func init() {
	flag.Uint64Var(&reqPerM, "ReqPMin", 600, "Max amount of requests")
	flag.Parse()
}

func main() {
	ExecuteRequests()
}

func HttpPostRequest() {
	result := make(chan Result)
	go func() {
		select {
		case <-time.After(TIMEOUT):
			log.Println("HttpPostRequest: timeout occurred, server could potentially lose message")
		case r := <-result:
			log.Println("HttpPostRequest: got response", r)
		}
	}()
	start := time.Now()
	resp, err := http.Post(URL, contentType, bytes.NewBuffer([]byte(MakeEvent())))
	if err != nil {
		//under heavy load it might happen
		//tcp: lookup localhost: device or resource busy
		log.Println("HttpPostRequest: Couldn't send a request to server.", err)
		result <- Result{true, 0, 0}
		return
	}
	result <- Result{false, resp.StatusCode, float64(time.Since(start).Seconds() * 1000)}
	defer resp.Body.Close()
}

func MakeEvent() string {
	val, err := uuid.NewV4()
	if err != nil {
		return `{"token":"invalid-token"}`
	}
	return fmt.Sprintf(`{"token":"%s"}`, val)
}

func ExecuteRequests() {
	numProcs := routines * runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numProcs)
	for p := 0; p < numProcs; p++ {
		go func() {
			defer wg.Done()
			//uniformly send requests
			limiter := time.Tick(time.Minute / time.Duration(reqPerM/uint64(numProcs)))
			for int64(atomic.AddUint64(&reqPerM, ^uint64(0))) >= 0 {
				<-limiter
				go HttpPostRequest()
			}
		}()
	}
	wg.Wait()
}
