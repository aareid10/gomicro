package micronum

import (
  "os"
	"context"
  "strconv"
	"encoding/json"
	"errors"
  "log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)




// NumberService provides operations on strings.
type NumberService interface {
	FibonacciCon(string) ([]int, error)
}


// TODO: Order response before sending

// numberService is a concrete implementation of StringService
type numberService struct{}

func (numberService) FibonacciCon(s string) ([]int, error) {
  layers, err  := strconv.Atoi(s)
  response := make([]int, layers)

  log.Print("layers", layers)

  if err == nil {
    jobs    := make(chan int, layers)
    results := make(chan int, layers)
    log.Print("setup done.")

    for w := 0; w < layers; w++ { go worker(jobs, results) }
    for i := 0; i < layers; i++ { jobs <- i }
    close(jobs)
    log.Print("jobs enqueue done.")

    for j := 0; j < layers; j++ { response[j] = <- results }
    log.Print("jobs results done.")

    return response, nil
  } else {
    log.Print("layer cast error.")
    return nil, errors.New("empty string")
  }

}

func worker(jobs <-chan int, results chan<- int)  {
  for n := range jobs { results <- fib(n) }
}

func fib(n int) int {
  if n <= 1 { return n }
  return fib(n - 1) + fib(n - 2)
}



// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")




// For each method, we define request and response structs
type fibonacciConRequest struct {
	F string `json:"f"`
}

type fibonacciConResponse struct {
	V   []int `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}




// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makefibonacciConEndpoint(svc NumberService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(fibonacciConRequest)
		v, err := svc.FibonacciCon(req.F)
		if err != nil {
			return fibonacciConResponse{v, err.Error()}, nil
		}
		return fibonacciConResponse{v, ""}, nil
	}
}




// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func Initialize() {
	svc := numberService{}

	fibonacciConHandler := httptransport.NewServer(
		makefibonacciConEndpoint(svc),
		decodefibonacciConRequest,
		encodeResponse,
	)

	http.Handle("/fibcon", fibonacciConHandler)

  port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Fatal(http.ListenAndServe(":"+ port, nil))
}

func decodefibonacciConRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request fibonacciConRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
