# 1million
Server, client may potentially consumes, produces 1 million requests per minute respectively (depends on your machine)
*keep in mind requests will be uniformly send within a minute
## Discussions on reddit
https://www.reddit.com/r/golang/comments/8torkz/1million_requests/
## Run tests within each directory (server/, client/)
- Unit tests
  -`go test -v`
- Benchmark tests
  -`go test -bench=.`

## How to build?
- Server
  - `go build -o server *.go`
- Client
  - `go build -o client *.go`

## How to run?
- `time ./server`
- `time ./client -ReqPMin 600` (600 requests within a minute)

### The benchmarks results
Running on *AMD A8-5557M Quad Core with turbo core technology up to 3.1 GHz*

#### Creating new workers benchmark test results:
![Image of New_workers benchmark tests](https://github.com/volodimyr/1million/blob/master/pictures/new_workerks.png)

#### Adding workers benchmark test results:
![Image of New_workers benchmark tests](https://github.com/volodimyr/1million/blob/master/pictures/add_workers.png)

#### Shutdown workers benchmark test results:
![Image of New_workers benchmark tests](https://github.com/volodimyr/1million/blob/master/pictures/shutdown_workers.png)
